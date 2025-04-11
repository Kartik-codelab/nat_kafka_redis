package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"nat_kafka_redis/domain"

	"github.com/go-redis/redis/v8"
)

type ClientRedis struct {
	client *redis.Client
}

type nat_data stuct {
	Latitude	float64 `json:"Latitude"`
	Longitude	float64 `json:"Longitude"`
	FixTime	time.time	`json:"Fixtime"`
}

func NewRedisClient(addr string, db int) *ClientRedis {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})

	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}
	return &ClientRedis{client: client}
}

func (r *ClientRedis) StorePosition(position domain.PositionPacket) error {
	ctx := context.Background()
	mobileCountryCode := position.Position.Network.CellTowers.MobileCountryCode
	mobileNetworkCode := position.Position.Network.CellTowers.MobileNetworkCode
	networkId := position.Position.Network.CellTowers.LocalAreaCode
	cellId := position.Position.Network.CellTowers.CellId
	key := fmt.Sprintf("bms_%v_%v_%v_%v", mobileCountryCode, mobileNetworkCode, networkId, cellId)

	data, err := json.Marshal(position.Position.Network)
	if err != nil {
		return err
	}

	nat_data := {
		Latitude: position.Position.Latitude,
		Longitude: position.Position.Longitude,
		FixTime: position.Position.FixTime
	}

	err = r.client.Set(ctx, key, nat_data, 0).Err()
	if err != nil {
		return err
	}
	log.Printf("Stored position to redis: %v", key, string(data))
	return nil
}

func (r *ClientRedis) Close() {
	r.client.Close()
}
