package kafka

import (
	"encoding/json"
	"log"
	"nat_kafka_redis/domain"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaProducer struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers[0]})
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v", ev.TopicPartition.Error)
				} else {
					log.Printf("Message delivered to %v", ev.TopicPartition)
				}
			}
		}
	}()

	return &KafkaProducer{
		producer: p,
		topic:    topic,
	}
}

func (p *KafkaProducer) SendPosition(position domain.PositionPacket) error {
	data, err := json.Marshal(position)
	if err != nil {
		return err
	}

	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          data,
	}, nil)
	if err != nil {
		return err
	}

	// Optionally flush to ensure delivery (for testing)
	p.producer.Flush(1000) // Wait up to 1 second
	log.Printf("Position packet sent: %s", string(data))
	return nil
}

func (p *KafkaProducer) Close() {
	p.producer.Close()
}
