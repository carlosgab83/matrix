package persisence

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
)

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
		return fmt.Errorf("failed to insert price %w in table %s: %w", price, pr.tableName(price.Symbol), err)
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
            id SERIAL PRIMARY KEY,
            price FLOAT NOT NULL,
            currency VARCHAR(10) NOT NULL,
            timestamp TIMESTAMP NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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

	hypertableQuery := fmt.Sprintf(`
    SELECT create_hypertable('%s', 'timestamp',
        if_not_exists => TRUE,
        migrate_data => TRUE
    )
`, tableName)

	_, err = pr.db.ExecContext(ctx, hypertableQuery)
	if err != nil {
		return fmt.Errorf("failed to create index on %s: %w", tableName, err)
	}

	return nil
}
