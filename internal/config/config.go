package config

import (
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv"
)

type Config struct {
	Env        string     `env:"APP_ENV" env-default:"local" env-required:"true"`
	DB         DBConfig   `env-prefix:"DB_"`
	HTTPServer HttpServer `env-prefix:"HTTP_"`
}

type DBConfig struct {
	Host     string `env:"HOST" env-required:"true"`
	Port     int    `env:"PORT" env-default:"5432"`
	User     string `env:"USER" env-default:"postgres"`
	Password string `env:"PASSWORD" env-required:"true"`
	Name     string `env:"NAME" env-required:"true"`
	SSLMode  string `env:"SSLMode" env-default:"disable"`
}

// return connection string for PostgreSQL
func (db DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.Name, db.SSLMode,
	)
}

type HttpServer struct {
	Address     string `env:"ADDRESS" env-default:"localhost:8080"`
	Timeout     int64  `env:"TIMEOUT" env-default:"4000000000"`      //nanoseconds
	IdleTimeout int64  `env:"IDLE_TIMEOUT" env-default:"6000000000"` //nanoseconds
}

func (h HttpServer) AsDuration() time.Duration {
	return time.Duration(h.Timeout)
}

func (h HttpServer) AsIdleDuration() time.Duration {
	return time.Duration(h.IdleTimeout)
}

func MustLoad() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("failed to load config from env: %s", err)
	}

	if cfg.Env == "" {
		log.Fatal("APP_ENV cannot be empty")
	}
	if cfg.DB.DSN() == "" {
		log.Fatal("database configuration is incomplete")
	}

	log.Printf("config loaded: env=%s, db=%s, http=%s",
		cfg.Env, cfg.DB.Host, cfg.HTTPServer.Address)

	return &cfg
}
