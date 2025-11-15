package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func NewFromEnv() (*sql.DB, error) {
	cfg := Config{
		Host:     getEnv("DB_HOST", "localhost"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "avito_review"),
	}

	portStr := getEnv("DB_PORT", "5432")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}
	cfg.Port = port

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}

	// Wait for the database to be ready
	const (
		maxAttempts = 10
		delay       = time.Second
	)

	for i := 1; i <= maxAttempts; i++ {
		if err := db.Ping(); err == nil {
			return db, nil
		}
		time.Sleep(delay)
	}

	_ = db.Close()
	return nil, fmt.Errorf("sql ping: database is not ready after %d attempts", maxAttempts)
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
