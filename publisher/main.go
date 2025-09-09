package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"test-task/internal/config"
	"test-task/internal/model"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/segmentio/kafka-go"
)

// topic, broker are configurable via env in config package

func main() {
	cntFakeData := config.PublisherCount()
	topic := config.PublisherTopic()
	broker := config.PublisherBroker()

	w := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	for i := 0; i < cntFakeData; i++ {
		var fake model.Order
		errF := gofakeit.Struct(&fake)
		if errF != nil {
			log.Fatal(errF)
		}

		data, err := json.Marshal(&fake)
		if err != nil {
			log.Fatalf("Failed to marshal fake order to JSON: %v", err)
		}

		if !json.Valid(data) {
			log.Fatalf("model.json contains invalid JSON")
		}
		fmt.Println(fake.OrderUID)

		err = w.WriteMessages(context.Background(),
			kafka.Message{
				Value: data,
			},
		)
		if err != nil {
			log.Fatalf("failed to write messages: %v", err)
		}
	}

	defer w.Close()

	log.Println("Message published to Kafka successfully!")
}
