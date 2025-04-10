// package usecase

// import (
// 	"nat_kafka_redis/domain"
// )

// type Consumer interface {
// 	ConsumerPositions(handler func(data domain.PositionPacket))
// }

// type consumerUsecase struct {
// 	kafkaConsumer Consumer
// }

// func NewConsumerUsecase(kafkaConsumer Consumer) Consumer {
// 	return &consumerUsecase{kafkaConsumer: kafkaConsumer}
// }

// func (c *consumerUsecase) ConsumerPositions(handler func(data domain.PositionPacket)) {
// 	c.kafkaConsumer.ConsumerPositions(handler)
// }
