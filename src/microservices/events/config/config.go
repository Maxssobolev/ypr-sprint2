package config

import "os"

type Config struct {
	Port         string
	KafkaBrokers string
}

func Load() *Config {
	return &Config{
		Port:         getEnv("PORT", "8082"),
		KafkaBrokers: getEnv("KAFKA_BROKERS", "kafka:9092"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
