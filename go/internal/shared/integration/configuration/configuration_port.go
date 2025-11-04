package configuration

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

func LoadConfig(cfg any, appName, path string) error {
	err := loadFileConfig(cfg, path)
	if err != nil {
		return fmt.Errorf("failed to load config for %s: %w", appName, err)
	}

	overlapEnvVars(cfg)

	return nil
}

func loadFileConfig(cfg any, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to parse JSON from %s: %w", path, err)
	}

	return nil
}

func overlapEnvVars(cfg any) {
	// Implementation to overlap environment variables onto cfg fields
	v := reflect.ValueOf(cfg).Elem()
	t := reflect.TypeOf(cfg).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		envVarKey := field.Tag.Get("env")
		if envVarKey != "" {
			if envValue := os.Getenv(envVarKey); envValue != "" {
				setFieldFromEnv(fieldValue, envValue)
			}
		}
	}
}

func setFieldFromEnv(fieldValue reflect.Value, envValue string) {
	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(envValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if valueInt, err := strconv.ParseInt(envValue, 10, 64); err == nil {
			fieldValue.SetInt(valueInt)
		}
	case reflect.Bool:
		if valueBool, err := strconv.ParseBool(envValue); err == nil {
			fieldValue.SetBool(valueBool)
		}
	}
}
