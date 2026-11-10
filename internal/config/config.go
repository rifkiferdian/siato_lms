package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName, Environment, Address, BaseURL, Timezone, SessionSecret, StoragePath string
	TrustedProxies                                                               []string
	SessionMaxAge                                                                int
	CookieSecure                                                                 bool
	DB                                                                           Database
}

type Database struct {
	Host, Port, Name, User, Password string
	MaxOpen, MaxIdle                 int
	ConnMaxLifetime                  time.Duration
}

func Load() (Config, error) {
	_ = godotenv.Load()
	c := Config{
		AppName: env("APP_NAME", "SIATO LMS"), Environment: env("APP_ENV", "development"),
		Address: env("APP_ADDR", ":8080"), BaseURL: env("APP_BASE_URL", "http://localhost:8080"),
		Timezone: env("APP_TIMEZONE", "Asia/Jakarta"), SessionSecret: os.Getenv("SESSION_SECRET"),
		StoragePath: env("STORAGE_PATH", "./storage/private"), SessionMaxAge: envInt("SESSION_MAX_AGE", 28800),
		CookieSecure: envBool("COOKIE_SECURE", false),
		DB:           Database{Host: env("DB_HOST", "127.0.0.1"), Port: env("DB_PORT", "3306"), Name: env("DB_NAME", "siato_lms"), User: env("DB_USER", "root"), Password: os.Getenv("DB_PASSWORD"), MaxOpen: envInt("DB_MAX_OPEN_CONNS", 25), MaxIdle: envInt("DB_MAX_IDLE_CONNS", 10), ConnMaxLifetime: envDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute)},
	}
	if v := strings.TrimSpace(os.Getenv("TRUSTED_PROXIES")); v != "" {
		c.TrustedProxies = strings.Split(v, ",")
	}
	if len(c.SessionSecret) < 32 {
		return Config{}, errors.New("SESSION_SECRET must contain at least 32 characters")
	}
	if _, err := time.LoadLocation(c.Timezone); err != nil {
		return Config{}, fmt.Errorf("invalid APP_TIMEZONE: %w", err)
	}
	return c, nil
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func envInt(k string, d int) int {
	v, e := strconv.Atoi(os.Getenv(k))
	if e == nil {
		return v
	}
	return d
}
func envBool(k string, d bool) bool {
	v, e := strconv.ParseBool(os.Getenv(k))
	if e == nil {
		return v
	}
	return d
}
func envDuration(k string, d time.Duration) time.Duration {
	v, e := time.ParseDuration(os.Getenv(k))
	if e == nil {
		return v
	}
	return d
}

func (d Database) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC&multiStatements=true", d.User, d.Password, d.Host, d.Port, d.Name)
}
