package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/brGuirra/greenlight/internal/data"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

const version = "1.0.0"

type Config struct {
	Environment          string `mapstructure:"ENVIRONMENT"`
	Port                 int    `mapstructure:"PORT"`
	DatabaseURL          string `mapstructure:"DATABASE_URL"`
	DatabaseMaxOpenConns int    `mapstructure:"DATABASE_MAX_OPEN_CONNECTIONS"`
	DatabaseMaxIdleConns int    `mapstructure:"DATABASE_MAX_IDLE_CONNECTIONS"`
	DatabaseMaxIdleTime  string `mapstructure:"DATABASE_MAX_IDLE_TIME"`
}

type application struct {
	logger *log.Logger
	config *Config
	models data.Models
}

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	cfg, err := loadConfig()
	if err != nil {
		logger.Fatal("Cannot load environent variables: ", err)
	}

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	logger.Print("database connection pool established")

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.Environment, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
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
