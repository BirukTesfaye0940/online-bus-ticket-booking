package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	// HTTP server
	HTTPPort string `mapstructure:"HTTP_PORT"`

	// Upstream gRPC service addresses
	AuthServiceAddr string `mapstructure:"AUTH_SERVICE_ADDR"`

	// Rate limiting
	RateLimitRequestsPerSecond float64 `mapstructure:"RATE_LIMIT_RPS"`
	RateLimitBurst             int     `mapstructure:"RATE_LIMIT_BURST"`
}

func Load() (*Config, error) {
	viper.SetDefault("HTTP_PORT", "8080")
	viper.SetDefault("AUTH_SERVICE_ADDR", "localhost:50051")
	viper.SetDefault("RATE_LIMIT_RPS", 10)
	viper.SetDefault("RATE_LIMIT_BURST", 20)

	// Load .env file if present (ignored in prod where real env vars are set)
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	_ = viper.ReadInConfig() // non-fatal if .env is missing

	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
