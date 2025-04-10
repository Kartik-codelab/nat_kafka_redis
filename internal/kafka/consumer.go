// package kafka

// import (
// 	"encoding/json"
// 	"log"
// 	"nat_kafka_redis/domain"

// 	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
// )

// type KafkaConsumer struct {
// 	consumer *kafka.Consumer
// 	topic    string
// }

// func NewKafkaConsumer(brokers []string, topic, groupId string) *KafkaConsumer {
// 	config := &kafka.ConfigMap{
// 		"bootstrap.servers": brokers[0],
// 		"auto.offset.reset": "earliest",
// 	}

// 	consumer, err := kafka.NewConsumer(config)
// 	if err != nil {
// 		log.Fatalf("Failed to create consumer: %v", err)
// 	}

// 	err = consumer.Subscribe(topic, nil)
// 	if err != nil {
// 		log.Fatalf("Failed to subscribe to topic %s: %v", topic, err)
// 	}

// 	return &KafkaConsumer{
// 		consumer: consumer,
// 		topic:    topic,
// 	}
// }

// func (c *KafkaConsumer) ConsumePosition(handler func(packet domain.PositionPacket)) {
// 	for {
// 		msg, err := c.consumer.ReadMessage(-1) // -1 means block indefinitely
// 		if err != nil {
// 			log.Printf("Consumer error: %v (%v)", err, msg)
// 			continue
// 		}

// 		var packet domain.PositionPacket
// 		if err := json.Unmarshal(msg.Value, &packet); err != nil {
// 			log.Printf("Failed to unmarshal position packet: %v", err)
// 			continue
// 		}

// 		handler(packet)
// 		log.Printf("Position packet received: %+v", packet)
// 	}
// }

// func (c *KafkaConsumer) Close() {
// 	c.consumer.Close()
// }
