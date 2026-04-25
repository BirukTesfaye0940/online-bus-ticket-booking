package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBURL    string `mapstructure:"DB_URL"`
	GRPCPort string `mapstructure:"GRPC_PORT"`
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("GRPC_PORT", "50052")

	// Load .env file if present
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	_ = viper.ReadInConfig()

	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
