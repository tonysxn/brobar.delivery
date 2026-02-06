package config

import (
	"github.com/tonysanin/brobar/pkg/helpers"
)

type Config struct {
	Port      string
	UploadDir string
}

func NewConfig() *Config {
	return &Config{
		Port:      helpers.GetEnv("SERVER_PORT", "3001"),
		UploadDir: helpers.GetEnv("UPLOAD_DIR", "./uploads"),
	}
}
