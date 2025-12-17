package config

import (
	"backend/core/constants"
	"backend/core/pkg/environment"
	"os"
)

type Config struct {
	Service          constants.Service
	Environment      *environment.Environment
	JWTSecret        string
	HTTPPort         string
	HTTPSchema       string
	LogLevel         string
	CORSOrigins      string
	PGHost           string
	PGPort           string
	PGDatabase       string
	PGUser           string
	PGPassword       string
	PGSSLMode        string
	RedisHost        string
	RedisPort        string
	RedisDatabase    string
	CentrifugoSchema string
	CentrifugoHost   string
	CentrifugoPort   string
	CentrifugoAPIKey string
	CentrifugoSecret string
}

func Load(service constants.Service) *Config {
	return &Config{
		Service:          service,
		Environment:      environment.New(getEnv("FIBER_ENVIRONMENT", "development")),
		JWTSecret:        getEnv("FIBER_JWT_SECRET", "secret"),
		HTTPPort:         getEnv("FIBER_HTTP_PORT", "3000"),
		HTTPSchema:       getEnv("FIBER_HTTP_SCHEMA", "http"),
		LogLevel:         getEnv("FIBER_LOG_LEVEL", "debug"),
		CORSOrigins:      getEnv("FIBER_CORS_ORIGINS", "http://localhost"),
		PGHost:           getEnv("POSTGRES_HOST", "localhost"),
		PGPort:           getEnv("POSTGRES_PORT", "5432"),
		PGDatabase:       getEnv("POSTGRES_DB", "postgres"),
		PGUser:           getEnv("POSTGRES_USER", "user"),
		PGPassword:       getEnv("POSTGRES_PASSWORD", "password"),
		PGSSLMode:        getEnv("POSTGRES_SSLMODE", "disable"),
		RedisHost:        getEnv("REDIS_HOST", "localhost"),
		RedisPort:        getEnv("REDIS_PORT", "6379"),
		RedisDatabase:    getEnv("REDIS_DB", "0"),
		CentrifugoSchema: getEnv("CENTRIFUGO_HTTP_SERVER_SCHEMA", "http"),
		CentrifugoHost:   getEnv("CENTRIFUGO_HTTP_SERVER_HOST", "localhost"),
		CentrifugoPort:   getEnv("CENTRIFUGO_HTTP_SERVER_PORT", "8000"),
		CentrifugoAPIKey: getEnv("CENTRIFUGO_HTTP_API_KEY", "fake_key"),
		CentrifugoSecret: getEnv("CENTRIFUGO_CLIENT_TOKEN_HMAC_SECRET_KEY", "fake_secret"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}
