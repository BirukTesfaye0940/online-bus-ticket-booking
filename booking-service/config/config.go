package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBURL              string `mapstructure:"DB_URL"`
	GRPCPort           string `mapstructure:"GRPC_PORT"`
	RedisURL           string `mapstructure:"REDIS_URL"`
	PaymentServiceAddr string `mapstructure:"PAYMENT_SERVICE_ADDR"`
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("GRPC_PORT", "50053")
	viper.SetDefault("REDIS_URL", "localhost:6379")
	viper.SetDefault("PAYMENT_SERVICE_ADDR", "localhost:50054")

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
