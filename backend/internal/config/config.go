// Package config
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	Env  string

	FrontEndURL string

	SupabaseURL        string
	SupabaseServiceKey string
	SupabaseJWTSecret  string

	DatabaseURL string
}

func LoadCfg() (*Config, error) {
	// if err := godotenv.Load(); err != nil {
	// 	return nil, err
	// }
	_ = godotenv.Load()

	cfg := &Config{
		Port: os.Getenv("PORT"),
		Env:  os.Getenv("ENV"),

		FrontEndURL: os.Getenv("FRONTEND_URL"),

		SupabaseURL:        os.Getenv("SUPABASE_URL"),
		SupabaseServiceKey: os.Getenv("SUPABASE_SERVICE_KEY"),
		SupabaseJWTSecret:  os.Getenv("SUPABASE_JWT_SECRET"),

		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Port == "" {
		return fmt.Errorf("PORT .env variable missing")
	}

	if c.Env == "" {
		return fmt.Errorf("ENV .env variable missing")
	}

	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL .env variable missing")
	}

	if c.SupabaseURL == "" {
		return fmt.Errorf("SUPABASE_URL .env variable missing")
	}

	if c.SupabaseServiceKey == "" {
		return fmt.Errorf("SUPABASE_SERVICE_KEY .env variable missing")
	}

	if c.FrontEndURL == "" {
		return fmt.Errorf("FRONTEND_URL .env variable missing")
	}

	if c.SupabaseJWTSecret == "" {
		return fmt.Errorf("SUPABASE_JWT_SECRET .env variable missing")
	}

	return nil
}
