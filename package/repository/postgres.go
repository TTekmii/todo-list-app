package repository

import (
	"fmt"
	// "log"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode,
	)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
// 	db, err := sqlx.Open("postgres", fmt.Sprintf(
// 		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
// 		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode,
// 	))
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = db.Ping()
// 	if err != nil {
// 		db.Close()
// 		return nil, err
// 	}

// 	return db, nil
// }

// func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
// 	return sqlx.Connect(
// 		"postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
// 		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode,
// 		))
// }
