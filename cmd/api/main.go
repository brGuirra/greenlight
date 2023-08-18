package main

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/brGuirra/greenlight/internal/data"
	"github.com/brGuirra/greenlight/internal/jsonlog"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

const version = "1.0.0"

type Config struct {
	Environment          string  `mapstructure:"ENVIRONMENT"`
	Port                 int     `mapstructure:"PORT"`
	DatabaseURL          string  `mapstructure:"DATABASE_URL"`
	DatabaseMaxOpenConns int     `mapstructure:"DATABASE_MAX_OPEN_CONNECTIONS"`
	DatabaseMaxIdleConns int     `mapstructure:"DATABASE_MAX_IDLE_CONNECTIONS"`
	DatabaseMaxIdleTime  string  `mapstructure:"DATABASE_MAX_IDLE_TIME"`
	RateLimitRPS         float64 `mapstructure:"RATE_LIMIT_RPS"`
	RateLimitBurst       int     `mapstructure:"RATE_LIMIT_BURST"`
	RatelimitEnabled     bool    `mapstructure:"RATE_LIMIT_ENABLED"`
}

type application struct {
	logger *jsonlog.Logger
	config *Config
	models data.Models
}

func main() {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	cfg, err := loadConfig()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func loadConfig() (*Config, error) {
	cfg := Config{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func openDB(cfg *Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DatabaseMaxOpenConns)
	db.SetMaxIdleConns(cfg.DatabaseMaxIdleConns)

	duration, err := time.ParseDuration(cfg.DatabaseMaxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
