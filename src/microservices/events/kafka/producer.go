package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers),
			Balancer: &kafka.LeastBytes{},
			Async:    false,
		},
	}
}

type MessageResult struct {
	Partition int
	Offset    int64
}

func (p *Producer) SendMessage(ctx context.Context, topic string, key []byte, value []byte) (MessageResult, error) {
	err := p.writer.WriteMessages(ctx,
		kafka.Message{
			Topic: topic,
			Key:   key,
			Value: value,
		},
	)

	if err != nil {
		return MessageResult{}, err
	}

	conn, err := kafka.DialLeader(context.Background(), "tcp", p.writer.Addr.String(), topic, 0)
	if err != nil {
		log.Printf("Failed to get Kafka metadata: %v", err)
		return MessageResult{Partition: 0, Offset: 0}, nil
	}
	defer conn.Close()

	offset, err := conn.ReadLastOffset()
	if err != nil {
		log.Printf("Failed to get last offset: %v", err)
		return MessageResult{Partition: 0, Offset: 0}, nil
	}

	return MessageResult{
		Partition: 0,
		Offset:    offset,
	}, nil
}

func (p *Producer) Close() {
	p.writer.Close()
}
