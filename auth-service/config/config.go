package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBURL          string        `mapstructure:"DB_URL"`
	GRPCPort       string        `mapstructure:"GRPC_PORT"`
	JWTSecret      string        `mapstructure:"JWT_SECRET"`
	TokenDuration  time.Duration `mapstructure:"TOKEN_DURATION"`
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("GRPC_PORT", "50051")
	viper.SetDefault("TOKEN_DURATION", 24*time.Hour)
	viper.SetDefault("JWT_SECRET", "super-secret-key")

	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
