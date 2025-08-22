package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/segmentio/kafka-go"
)

const (
	topic  = "orders"
	broker = "localhost:9092"
)

func main() {
	// Чтение модели из файла
	data, err := ioutil.ReadFile("model.json")
	if err != nil {
		log.Fatalf("Failed to read model file: %v", err)
	}

	// Проверка валидности JSON
	if !json.Valid(data) {
		log.Fatalf("model.json contains invalid JSON")
	}

	// Настройка писателя Kafka
	w := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	defer w.Close()

	// Отправка сообщения
	err = w.WriteMessages(context.Background(),
		kafka.Message{
			Value: data,
		},
	)
	if err != nil {
		log.Fatalf("failed to write messages: %v", err)
	}

	log.Println("Message published to Kafka successfully!")
}
