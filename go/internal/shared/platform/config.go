package platform

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadConfig loads JSON from a file and unmarshals it into the provided struct
// cfg should be a pointer to the struct you want to fill
// Example: LoadConfig(&myConfig, "config.json")
func LoadConfig(cfg any, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to parse JSON from %s: %w", path, err)
	}

	return nil
}
