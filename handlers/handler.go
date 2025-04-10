package handler

import (
	"log"
	"nat_kafka_redis/domain"
	"nat_kafka_redis/redis"
)

type PositionHandler struct {
	redisClient *redis.ClientRedis
}

func NewPositionHandler(redisClient *redis.ClientRedis) *PositionHandler {
	return &PositionHandler{redisClient: redisClient}
}

func (h *PositionHandler) Handle(packet domain.PositionPacket) {
	log.Printf("PositionHandler::Handle(%+v)\n", packet)

	err := h.redisClient.StorePosition(packet)
	if err != nil {
		log.Printf("Failed to store position packet in Redis: %v", err)
	}
}
