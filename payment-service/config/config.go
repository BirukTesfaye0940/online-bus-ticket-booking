package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBURL              string `mapstructure:"DB_URL"`
	GRPCPort           string `mapstructure:"GRPC_PORT"`
	StripeSecretKey    string `mapstructure:"STRIPE_SECRET_KEY"`
	StripeWebhookSecret string `mapstructure:"STRIPE_WEBHOOK_SECRET"`
	BookingServiceAddr string `mapstructure:"BOOKING_SERVICE_ADDR"`
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("GRPC_PORT", "50054")
	viper.SetDefault("BOOKING_SERVICE_ADDR", "localhost:50053")

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
