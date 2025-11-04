package domain

type CommonConfig struct {
	LogLevel    string `json:"log_level" env:"MATRIX_LOG_LEVEL"`
	LogFilePath string `json:"log_file_path" env:"MATRIX_LOG_FILE_PATH"`
}
