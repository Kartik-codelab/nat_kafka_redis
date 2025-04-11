package handler

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"nat_kafka_redis/internal/domain"
	"nat_kafka_redis/internal/infrastructure/config"
	kafkainfra "nat_kafka_redis/internal/infrastructure/kafka"
	redisinfra "nat_kafka_redis/internal/infrastructure/redis"
	"nat_kafka_redis/internal/usecase"
)

func TestPositionHandlerIntegration(t *testing.T) {
	log.SetOutput(os.Stdout)

	ctx := context.Background()

	// Start Kafka container
	kafkaContainer, err := kafka.Run(ctx, "confluentinc/cp-kafka:latest",
		kafka.WithClusterID("test-cluster"),
		kafka.WithZookeeper(),
		testcontainers.WithEnv(map[string]string{
			"KAFKA_ADVERTISED_LISTENERS":             "PLAINTEXT://localhost:9092",
			"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR": "1",
		}),
	)
	assert.NoError(t, err)
	defer kafkaContainer.Terminate(ctx)

	kafkaBroker, err := kafkaContainer.Brokers(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, kafkaBroker)

	// Start Redis container
	redisContainer, err := redis.Run(ctx, "redis:latest",
		redis.WithPassword("redis_password"),
	)
	assert.NoError(t, err)
	defer redisContainer.Terminate(ctx)

	redisAddr, err := redisContainer.Endpoint(ctx, "")
	assert.NoError(t, err)

	// Configure services
	cfg := &config.Config{
		KafkaBrokers:  kafkaBroker[0],
		KafkaTopic:    "iot-position",
		KafkaGroupID:  "test-consumer-group",
		RedisAddr:     redisAddr,
		RedisPassword: "redis_password",
		RedisDB:       0,
	}

	producer := kafkainfra.NewKafkaProducer([]string{cfg.KafkaBrokers}, cfg.KafkaTopic)
	defer producer.Close()

	consumer := kafkainfra.NewKafkaConsumer([]string{cfg.KafkaBrokers}, cfg.KafkaTopic, cfg.KafkaGroupID)
	defer consumer.Close()

	redisClient := redisinfra.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	defer redisClient.Close()

	handler := NewPositionHandler(redisClient)
	consumerUsecase := usecase.NewConsumerUsecase(consumer)
	producerUsecase := usecase.NewProducerUsecase(producer)

	// Start consumer in a goroutine
	done := make(chan bool)
	go func() {
		consumerUsecase.ConsumePositionPackets(handler.Handle)
		done <- true
	}()

	// Wait for consumer to start
	time.Sleep(2 * time.Second)

	// Send a test packet
	now := time.Now()
	typeStr := "GPS"
	testPacket := domain.PositionPacket{
		Position: domain.PositionData{
			ID:         1,
			Attributes: map[string]interface{}{"battery": 75},
			DeviceID:   1,
			Type:       &typeStr,
			Protocol:   "osmand",
			ServerTime: now,
			DeviceTime: now.Add(-time.Second),
			FixTime:    now.Add(-2 * time.Second),
			Outdated:   false,
			Valid:      true,
			Latitude:   37.7749,
			Longitude:  -122.4194,
			Altitude:   10.5,
			Speed:      5.0,
			Course:     180.0,
			Address:    nil,
			Accuracy:   1.0,
			Network: domain.NetworkData{
				RadioType:  "LTE",
				ConsiderIP: false,
				CellTowers: []domain.CellTower{
					{CellID: 123, LocationAreaCode: 456, MobileCountryCode: 310, MobileNetworkCode: 410},
				},
			},
		},
		Device: domain.DeviceData{
			ID:         1,
			Attributes: map[string]interface{}{"firmware": "v1.0"},
			GroupID:    0,
			Name:       "Tracker-001",
			UniqueID:   "iot-001",
			Status:     "online",
			LastUpdate: now,
			PositionID: 1,
			GeofenceIds: []int64{1, 2},
			Phone:      "123-456-7890",
			Model:      "GT-100",
			Contact:    "admin@example.com",
			Category:   "vehicle",
			Disabled:   false,
		},
	}
	err = producerUsecase.SendPositionPacket(testPacket)
	assert.NoError(t, err)

	// Wait for consumption and Redis storage
	time.Sleep(2 * time.Second)

	// Verify device key
	ctxRedis := context.Background()
	deviceKey := "device_" + testPacket.Device.UniqueID
	redisData, err := redisClient.client.Get(ctxRedis, deviceKey).Result()
	assert.NoError(t, err)

	var storedPacket domain.PositionPacket
	err = json.Unmarshal([]byte(redisData), &storedPacket)
	assert.NoError(t, err)

	assert.Equal(t, testPacket.Device.UniqueID, storedPacket.Device.UniqueID, "Device UniqueID should match")
	assert.Equal(t, testPacket.Position.Latitude, storedPacket.Position.Latitude, "Device Latitude should match")
	assert.Equal(t, testPacket.Position.Longitude, storedPacket.Position.Longitude, "Device Longitude should match")
	assert.Equal(t, testPacket.Position.FixTime, storedPacket.Position.FixTime, "Device FixTime should match")

	// Verify cell tower key
	cellTowerKey := "cellTower_123_310_456_410"
	cellTowerData, err := redisClient.client.Get(ctxRedis, cellTowerKey).Result()
	assert.NoError(t, err)

	var storedCellTower struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		FixTime   string  `json:"fixTime"`
	}
	err = json.Unmarshal([]byte(cellTowerData), &storedCellTower)
	assert.NoError(t, err)

	expectedFixTime := testPacket.Position.FixTime.Format("2006-01-02T15:04:05Z07:00")
	assert.Equal(t, testPacket.Position.Latitude, storedCellTower.Latitude, "Cell tower Latitude should match")
	assert.Equal(t, testPacket.Position.Longitude, storedCellTower.Longitude, "Cell tower Longitude should match")
	assert.Equal(t, expectedFixTime, storedCellTower.FixTime, "Cell tower FixTime should match")

	// Optionally stop consumer to clean up
	consumer.Close()
}