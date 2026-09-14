package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	JWTExpiry   time.Duration
	Port        string
	UploadDir   string
	BcryptCost  int
	Migrations  string
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	secret := os.Getenv("JWT_SECRET")
	if len(secret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 bytes")
	}
	expiry := 24 * time.Hour
	if s := os.Getenv("JWT_EXPIRY"); s != "" {
		d, err := time.ParseDuration(s)
		if err != nil {
			return nil, fmt.Errorf("JWT_EXPIRY: %w", err)
		}
		expiry = d
	}
	cost := 12
	if s := os.Getenv("BCRYPT_COST"); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("BCRYPT_COST: %w", err)
		}
		cost = n
	}
	if cost < 12 {
		return nil, fmt.Errorf("BCRYPT_COST must be >= 12")
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	upload := os.Getenv("UPLOAD_DIR")
	if upload == "" {
		upload = "./data/uploads"
	}
	mig := os.Getenv("MIGRATIONS_DIR")
	if mig == "" {
		mig = "migrations"
	}
	return &Config{
		DatabaseURL: dsn,
		JWTSecret:   secret,
		JWTExpiry:   expiry,
		Port:        port,
		UploadDir:   upload,
		BcryptCost:  cost,
		Migrations:  mig,
	}, nil
}
