package config

import (
	"strconv"
	"strings"
)

const (
	defaultPort = 2810
	dbPort      = 5432
)

type Config struct {
	RootPath string
	App      App
	DB       DB
	Logs     Logs
}

type App struct {
	Name  string
	Env   string
	Port  int
	Host  string
	Debug bool
}

type DB struct {
	Connection string
	Host       string
	Port       int
	Username   string
	Password   string
	Database   string
}

type Logs struct {
	Level string
}

func NewConfig(lookenv func(string) (string, bool)) Config {
	getString := getEnvAsScalar[string](lookenv)
	getInt := getEnvAsScalar[int](lookenv)
	getBool := getEnvAsScalar[bool](lookenv)

	return Config{
		App: App{
			Name:  getString("APP_NAME", "golangapp"),
			Env:   getString("APP_ENV", "local"),
			Port:  getInt("APP_PORT", defaultPort),
			Host:  getString("APP_HOST", "localhost"),
			Debug: getBool("APP_DEBUG", false),
		},
		DB: DB{
			Connection: getString("DB_CONNECTION", "postgres"),
			Host:       getString("DB_HOST", "localhost"),
			Port:       getInt("DB_PORT", dbPort),
			Username:   getString("DB_USERNAME", "postgres"),
			Password:   getString("DB_PASSWORD", "password"),
			Database:   getString("DB_DATABASE", "postgres"),
		},
		Logs: Logs{
			// Level: nope, debug, info, warn, error
			Level: getString("LOG_LEVEL", "info"),
		},
	}
}

func (c Config) IsProduction() bool {
	return strings.Contains(strings.ToLower(c.App.Env), "prod")
}

func (c Config) IsDevelopment() bool {
	return c.App.Env == "local" || c.App.Env == "dev"
}

func (c Config) IsTest() bool {
	return strings.Contains(strings.ToLower(c.App.Env), "test")
}

func getEnvAsScalar[T any](lookenv func(string) (string, bool)) func(string, T) T {
	return func(key string, fallback T) T {
		value, ok := lookenv(key)

		if !ok || value == "" {
			return fallback
		}

		var result any

		switch any(fallback).(type) {
		case string:
			result = value
		case int:
			valueToInt, err := strconv.Atoi(value)
			if err != nil {
				return fallback
			}

			result = valueToInt
		case bool:
			valueToBool, err := strconv.ParseBool(value)
			if err != nil {
				return fallback
			}

			result = valueToBool
		}

		return result.(T)
	}
}
