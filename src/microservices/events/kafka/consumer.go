package kafka

import (
	"context"
	"encoding/json"
	"events/models"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  groupID,
			Topic:    "movie-events,user-events,payment-events",
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),
	}
}

func (c *Consumer) Consume(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping Kafka consumer...")
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			switch msg.Topic {
			case "movie-events":
				var event models.EventDetail
				if err := json.Unmarshal(msg.Value, &event); err != nil {
					log.Printf("Error decoding movie event: %v", err)
					continue
				}
			case "user-events":
				var event models.EventDetail
				if err := json.Unmarshal(msg.Value, &event); err != nil {
					log.Printf("Error decoding user event: %v", err)
					continue
				}
			case "payment-events":
				var event models.EventDetail
				if err := json.Unmarshal(msg.Value, &event); err != nil {
					log.Printf("Error decoding user event: %v", err)
					continue
				}
			}
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
