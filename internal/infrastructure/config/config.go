package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port           string
	Env            string
	AllowedOrigins []string
	BootstrapAdmin BootstrapAdminConfig
	DB             DBConfig
	JWT            JWTConfig
}

type DBConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type JWTConfig struct {
	Secret string
	TTL    time.Duration
}

type BootstrapAdminConfig struct {
	Email    string
	Password string
	Name     string
}

func Load() Config {
	loadDotEnv(".env")
	return Config{
		Port:           env("APP_PORT", "8080"),
		Env:            env("APP_ENV", "development"),
		AllowedOrigins: envCSV("APP_ALLOWED_ORIGINS", "http://localhost:5184,http://127.0.0.1:5184"),
		BootstrapAdmin: BootstrapAdminConfig{
			Email:    env("BOOTSTRAP_ADMIN_EMAIL", "admin@dachcourier.co.uk"),
			Password: env("BOOTSTRAP_ADMIN_PASSWORD", "ChangeMe123!"),
			Name:     env("BOOTSTRAP_ADMIN_NAME", "Dach Super Admin"),
		},
		DB: DBConfig{
			Host:            env("DB_HOST", "localhost"),
			Port:            env("DB_PORT", "5432"),
			User:            env("DB_USER", "postgres"),
			Password:        env("DB_PASSWORD", "postgres"),
			Name:            env("DB_NAME", "dach"),
			SSLMode:         env("DB_SSLMODE", "disable"),
			MaxOpenConns:    envInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    envInt("DB_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime: envDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute),
			ConnMaxIdleTime: envDuration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute),
		},
		JWT: JWTConfig{
			Secret: env("JWT_SECRET", "change-me"),
			TTL:    envDuration("JWT_TTL", 24*time.Hour),
		},
	}
}

func loadDotEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if key == "" || os.Getenv(key) != "" {
			continue
		}
		_ = os.Setenv(key, value)
	}
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}

func env(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func envDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func envCSV(key string, fallback string) []string {
	raw := env(key, fallback)
	values := strings.Split(raw, ",")
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			out = append(out, value)
		}
	}
	return out
}
