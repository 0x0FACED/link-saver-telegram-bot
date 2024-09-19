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
	Link LinkServiceConfig
	PDF  PDFServiceConfig
}

type LinkServiceConfig struct {
	Host string
	Port string
}

type PDFServiceConfig struct {
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
			Link: LinkServiceConfig{
				Host: os.Getenv("GRPC_HOST"),
				Port: os.Getenv("GRPC_PORT"),
			},
			PDF: PDFServiceConfig{
				Host: os.Getenv("PDF_HOST"),
				Port: os.Getenv("PDF_PORT"),
			},
		},
		Telegram: TelegramConfig{
			Token: os.Getenv("BOT_TOKEN"),
		},
	}, nil
}
