package main

import (
	"log"
	"math/rand"
	"nat_kafka_redis/config"
	"nat_kafka_redis/domain"
	"nat_kafka_redis/internal/kafka"
	"nat_kafka_redis/internal/usecase"
	"time"
)

func main() {
	cfg := config.LoadConfig()

	producer := kafka.NewKafkaProducer([]string{cfg.KafkaBrokers}, cfg.KafkaTopic)
	defer producer.Close()

	producerUsecase := usecase.NewProducerUsecase(producer)

	for {
		now := time.Now()
		typeStr := "GPS"
		packet := domain.PositionPacket{
			Position: domain.PositionData{
				ID:         rand.Int63n(1000),
				Attributes: map[string]interface{}{"battery": 75},
				DeviceID:   1,
				Type:       &typeStr,
				Protocol:   "teltonika",
				ServerTime: now,
				DeviceTime: now.Add(-time.Second),
				FixTime:    now.Add(-2 * time.Second),
				Outdated:   false,
				Valid:      true,
				Latitude:   37.7749 + rand.Float64()*0.01,
				Longitude:  -122.4194 + rand.Float64()*0.01,
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
				ID:          1,
				Attributes:  map[string]interface{}{"firmware": "v1.0"},
				GroupID:     0,
				Name:        "Tracker-001",
				UniqueID:    "iot-001",
				Status:      "online",
				LastUpdate:  now,
				PositionID:  rand.Int63n(1000),
				GeofenceIds: []int64{1, 2},
				Phone:       "123-456-7890",
				Model:       "GT-100",
				Contact:     "admin@example.com",
				Category:    "vehicle",
				Disabled:    false,
			},
		}
		err := producerUsecase.SendPosition(packet)
		if err != nil {
			log.Println("Failed to send position:", err)
		}
		time.Sleep(5 * time.Second)
	}
}
