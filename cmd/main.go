package main

import (
	"log"
	"test-task/internal/cache"
	"test-task/internal/db"
	"test-task/internal/repository"
	"test-task/internal/service"
)

func main() {
	database := db.InitDB()
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error close Database: %v", err)
		}
	}()
	cache := cache.InitCache(database)
	service.NewOrderService(cache, repository.NewOrderRepository(database))
	log.Println("Service started")
}
