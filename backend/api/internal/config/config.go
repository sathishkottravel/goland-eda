// Package config loads service configuration from environment variables.
package config

import "os"

type Config struct {
	Port         string
	KafkaBrokers string
}

func Load() Config {
	return Config{
		Port:         getenv("PORT", "8080"),
		KafkaBrokers: getenv("KAFKA_BROKERS", "localhost:9092"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
