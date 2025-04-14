package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPServer HTTPServer
	PGConfig   PGConfig
	JWTSecret  string
	PassSecret string
}

type HTTPServer struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type PGConfig struct {
	DSN string
}

func New(path string) (*Config, error) {
	env := os.Getenv("ENV")
	if env != "prod" {
		err := godotenv.Load(path)
		if err != nil {
			return nil, err
		}
	}

	return &Config{
		HTTPServer: HTTPServer{
			Port:         os.Getenv("HTTP_PORT"),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		PGConfig: PGConfig{
			DSN: os.Getenv("PG_DSN"),
		},
		JWTSecret:  os.Getenv("JWT_SECRET"),
		PassSecret: os.Getenv("PASSWORD_SECRET"),
	}, nil
}
