package usecase

import "nat_kafka_redis/domain"

type Producer interface {
	SendPosition(position domain.PositionPacket) error
}

type producerUsecase struct {
	kafkaProducer Producer
}

func NewProducerUsecase(kafkaProducer Producer) Producer {
	return &producerUsecase{kafkaProducer: kafkaProducer}
}

func (p *producerUsecase) SendPosition(position domain.PositionPacket) error {
	return p.kafkaProducer.SendPosition(position)
}
