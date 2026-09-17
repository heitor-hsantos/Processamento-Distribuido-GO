package config

import "os"

type Config struct {
	AppName     string
	Port        int
	PostgresURL string
	SQSQueueURL string
	IssuerURL   string
	Audience    string
}

func Load() Config {
	return Config{
		AppName:     getEnv("APP_NAME", "jungle-gaming"),
		Port:        8080,
		PostgresURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/jungle_gaming"),
		SQSQueueURL: getEnv("SQS_QUEUE_URL", "http://localhost:4566/000000000000/betting-events"),
		IssuerURL:   getEnv("OIDC_ISSUER_URL", "http://localhost:8081/realms/jungle"),
		Audience:    getEnv("OIDC_AUDIENCE", "jungle-gaming"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
