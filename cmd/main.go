package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test-task/internal/cache"
	"test-task/internal/db"
	"test-task/internal/handlers"
	"test-task/internal/kafka"
	"test-task/internal/repository"
	"test-task/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	database := db.InitDB()
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error close Database: %v", err)
		}
	}()

	orderCache := cache.InitCache(database)
	orderRepo := repository.NewOrderRepository(database)
	orderService := service.NewOrderService(orderCache, orderRepo)

	ctx, cancel := context.WithCancel(context.Background())
	kafkaSubscriber := kafka.NewKafkaSubscriber(orderService)
	go kafkaSubscriber.Subscribe(ctx)
	defer kafkaSubscriber.Close()

	orderHandler := handlers.NewOrderHandler(orderService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/order/{order_uid}", orderHandler.GetOrder)

	fs := http.FileServer(http.Dir("./web"))
	r.Handle("/*", fs)

	log.Println("Service started on :8081")
	go func() {
		if err := http.ListenAndServe(":8081", r); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	cancel()
}
