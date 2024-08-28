package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPC     GRPCConfig
	Telegram TelegramConfig
}

type GRPCConfig struct {
	Host string
	Port string
}

type TelegramConfig struct {
	Token string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		GRPC: GRPCConfig{
			Host: os.Getenv("GRPC_HOST"),
			Port: os.Getenv("GRPC_PORT"),
		},
		Telegram: TelegramConfig{
			Token: os.Getenv("BOT_TOKEN"),
		},
	}, nil
}
