package main

import (
	"nat_kafka_redis/config"
	handler "nat_kafka_redis/handlers"
	"nat_kafka_redis/internal/kafka"
	"nat_kafka_redis/internal/usecase"
	"nat_kafka_redis/redis"
)

func main() {
	cfg := config.LoadConfig()

	redisClient := redis.NewRedisClient(cfg.RedisAddr, cfg.RedisDB)
	defer redisClient.Close()

	consumer := kafka.NewKafkaConsumer([]string{cfg.KafkaBrokers}, cfg.KafkaTopic, cfg.KafkaGroupId)
	defer consumer.Close()

	consumerUsecase := usecase.NewConsumerUsecase(consumer)
	positionHandler := handler.NewPositionHandler(redisClient)

	consumerUsecase.ConsumerPositions(positionHandler.Handle)
}
