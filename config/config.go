package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	KafkaBrokers string
	KafkaTopic   string
	KafkaGroupId string
	RedisAddr    string
	RedisDB      int
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		KafkaBrokers: os.Getenv("KAFKA_BROKERS"),
		KafkaTopic:   os.Getenv("KAFKA_TOPIC"),
		KafkaGroupId: os.Getenv("KAFKA_GROUP_ID"),
		RedisAddr:    os.Getenv("REDIS_ADDR"),
		RedisDB:      0,
	}
}
