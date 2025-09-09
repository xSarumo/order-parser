package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test-task/internal/cache"
	"test-task/internal/config"
	"test-task/internal/db"
	"test-task/internal/handlers"
	"test-task/internal/kafka"
	"test-task/internal/repository"
	"test-task/internal/service"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   config.CORSAllowedOrigins(),
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.Logger)
	r.Get("/order/{order_uid}", orderHandler.GetOrder)

	fs := http.FileServer(http.Dir(config.StaticDir()))
	r.Handle("/*", fs)

	srv := &http.Server{
		Addr:    config.HTTPAddr(),
		Handler: r,
	}

	go func() {
		log.Println("Service started on", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("could not listen on %s: %v\n", srv.Addr, err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Println("Server exited properly")
}
