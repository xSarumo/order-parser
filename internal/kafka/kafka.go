package kafka

import (
	"context"
	"encoding/json"
	"log"
	"test-task/internal/config"
	"test-task/internal/model"
	"test-task/internal/service"

	"github.com/segmentio/kafka-go"
)

type KafkaSubscriber struct {
	reader  *kafka.Reader
	service *service.OrderService
}

func NewKafkaSubscriber(service *service.OrderService) *KafkaSubscriber {
	broker := config.KafkaBroker()
	topic := config.KafkaTopic()
	groupID := config.KafkaGroupID()
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		GroupID: groupID,
		Topic:   topic,
	})

	return &KafkaSubscriber{
		reader:  r,
		service: service,
	}
}

func (ks *KafkaSubscriber) Subscribe(ctx context.Context) {
	log.Println("Subscribed to Kafka topic:", config.KafkaTopic())
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping Kafka subscriber...")
			return
		default:
			m, err := ks.reader.FetchMessage(ctx)
			if err != nil {
				log.Printf("could not fetch message: %v", err)
				continue
			}

			var order model.Order
			if err := json.Unmarshal(m.Value, &order); err != nil {
				log.Printf("Failed to unmarshal order: %v", err)
				continue
			}

			ks.service.ProcessNewOrder(order)

			if err := ks.reader.CommitMessages(ctx, m); err != nil {
				log.Printf("failed to commit messages: %v", err)
			}
		}
	}
}

func (ks *KafkaSubscriber) Close() {
	if ks.reader != nil {
		if err := ks.reader.Close(); err != nil {
			log.Printf("failed to close kafka reader: %v", err)
		}
	}
}
