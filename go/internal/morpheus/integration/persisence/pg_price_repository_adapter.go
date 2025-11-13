package persisence

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
)

type Scale struct {
	suffix           string
	interval         string
	scheduleInterval string
	startOffset      string
	endOffset        string
}

type PgPriceRepository struct {
	db *sql.DB
}

func NewPgPriceRepository(connStr string) (*PgPriceRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db.SetMaxOpenConns(25)                 // Maximum number of open connections
	db.SetMaxIdleConns(5)                  // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Maximum lifetime of a connection

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PgPriceRepository{
		db: db,
	}, nil
}

func (pr *PgPriceRepository) InsertPrice(ctx context.Context, price shared_domain.Price) error {
	if err := pr.ensureTableExists(ctx, pr.tableName(price.Symbol)); err != nil {
		return fmt.Errorf("failed to create %s table: %w", pr.tableName(price.Symbol), err)
	}

	insertQuery := fmt.Sprintf(`
		INSERT INTO %s (price, currency, timestamp)
		VALUES ($1, $2, $3)
	`, pr.tableName(price.Symbol))

	_, err := pr.db.ExecContext(ctx, insertQuery, price.Price, price.Currency, price.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to insert price %v in table %s: %w", price, pr.tableName(price.Symbol), err)
	}

	return nil
}

func (pr *PgPriceRepository) Close() error {
	return pr.db.Close()
}

func (pr *PgPriceRepository) tableName(symbol string) string {
	return strings.ToLower(symbol)
}

func (pr *PgPriceRepository) ensureTableExists(ctx context.Context, tableName string) error {
	checkQuery := `
        SELECT EXISTS (
            SELECT FROM information_schema.tables
            WHERE table_schema = 'public'
            AND table_name = $1
        )
    `

	var exists bool
	err := pr.db.QueryRowContext(ctx, checkQuery, tableName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if table exists: %w", err)
	}

	if exists {
		return nil
	}

	createQuery := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            price FLOAT NOT NULL,
            currency text NOT NULL,
            timestamp TIMESTAMPTZ NOT NULL,
            created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
        )
    `, tableName)

	_, err = pr.db.ExecContext(ctx, createQuery)
	if err != nil {
		return fmt.Errorf("failed to create table %s: %w", tableName, err)
	}

	indexQuery := fmt.Sprintf(`
        CREATE INDEX IF NOT EXISTS idx_%s_timestamp
        ON %s (timestamp DESC)
    `, tableName, tableName)

	_, err = pr.db.ExecContext(ctx, indexQuery)
	if err != nil {
		return fmt.Errorf("failed to create index on %s: %w", tableName, err)
	}

	err = pr.createHypertableSetup(ctx, tableName)
	if err != nil {
		return fmt.Errorf("failed to create hypertable setup on %s: %w", tableName, err)
	}

	return nil
}

func (pr *PgPriceRepository) createHypertableSetup(ctx context.Context, tableName string) error {
	err := pr.createAndSetupHypertable(ctx, tableName)
	if err != nil {
		return fmt.Errorf("failed to add refresh policy on %s_5m: %w", tableName, err)
	}

	scales := []Scale{
		{
			suffix:           "1m",
			interval:         "1 minute",
			scheduleInterval: "10 second",
			startOffset:      "10 minute",
			endOffset:        "0 minute",
		},
		{
			suffix:           "5m",
			interval:         "5 minute",
			scheduleInterval: "10 second",
			startOffset:      "50 minute",
			endOffset:        "0 minute",
		},
	}

	err = pr.createAndsetupMaterializedView(ctx, tableName, scales)
	if err != nil {
		return fmt.Errorf("failed to add refresh policy on %s_5m: %w", tableName, err)
	}

	return nil
}

func (pr *PgPriceRepository) createAndSetupHypertable(ctx context.Context, tableName string) error {
	hypertableQuery := fmt.Sprintf(`
    SELECT create_hypertable('%s', 'timestamp',
        if_not_exists => TRUE,
        migrate_data => TRUE
    )
`, tableName)

	_, err := pr.db.ExecContext(ctx, hypertableQuery)
	if err != nil {
		return fmt.Errorf("failed to create index on %s: %w", tableName, err)
	}

	// Add retention policy - keep data for 365 days
	retentionQuery := fmt.Sprintf(`
    SELECT add_retention_policy('%s', INTERVAL '365 days')
`, tableName)

	_, err = pr.db.ExecContext(ctx, retentionQuery)
	if err != nil {
		return fmt.Errorf("failed to add retention policy on %s: %w", tableName, err)
	}

	// Enable compression on the table
	compressionQuery := fmt.Sprintf(`
    ALTER TABLE %s SET (
        timescaledb.compress,
        timescaledb.compress_segmentby = 'currency',
        timescaledb.compress_orderby = 'timestamp DESC'
    )
`, tableName)

	_, err = pr.db.ExecContext(ctx, compressionQuery)
	if err != nil {
		return fmt.Errorf("failed to enable compression on %s: %w", tableName, err)
	}

	// Add compression policy - compress chunks older than 6 months
	// Only enables compression of tableName.
	// Enable on materialized views too if needed.
	compressionPolicyQuery := fmt.Sprintf(`
    SELECT add_compression_policy('%s', INTERVAL '6 months')
`, tableName)

	_, err = pr.db.ExecContext(ctx, compressionPolicyQuery)
	if err != nil {
		return fmt.Errorf("failed to add compression policy on %s: %w", tableName, err)
	}

	return nil
}

func (pr *PgPriceRepository) createAndsetupMaterializedView(ctx context.Context, tableName string, scales []Scale) error {
	for _, scale := range scales {
		// Create continuous aggregate for 1-minute intervals
		aggregateQuery := fmt.Sprintf(`
    CREATE MATERIALIZED VIEW IF NOT EXISTS %s_%s
    WITH (timescaledb.continuous) AS
    SELECT
        time_bucket('%s', timestamp) AS bucket,
        AVG(price) AS avg_price,
        MAX(price) AS max_price,
        MIN(price) AS min_price,
        FIRST(price, timestamp) AS open_price,
        LAST(price, timestamp) AS close_price,
        COUNT(*) AS prices
    FROM %s
    GROUP BY bucket
`, tableName, scale.suffix, scale.interval, tableName)

		_, err := pr.db.ExecContext(ctx, aggregateQuery)
		if err != nil {
			return fmt.Errorf("failed to create continuous aggregate on %s: %w", tableName, err)
		}

		refreshPolicyQuery := fmt.Sprintf(`
    SELECT add_continuous_aggregate_policy('%s_%s',
        start_offset => INTERVAL '%s',
        end_offset => INTERVAL '%s',
        schedule_interval => INTERVAL '%s')
`, tableName, scale.suffix, scale.startOffset, scale.endOffset, scale.scheduleInterval)

		_, err = pr.db.ExecContext(ctx, refreshPolicyQuery)
		if err != nil {
			return fmt.Errorf("failed to add refresh policy on %s_%s: %w", tableName, scale.suffix, err)
		}
	}

	return nil
}
