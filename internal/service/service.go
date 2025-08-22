package service

import (
	"errors"
	"log"
	"test-task/internal/cache"
	"test-task/internal/model"
	"test-task/internal/repository"
)

type OrderService struct {
	cache *cache.LRU_Cache
	repo  *repository.OrderRepository
}

func NewOrderService(cache *cache.LRU_Cache, repo *repository.OrderRepository) *OrderService {
	return &OrderService{
		cache: cache,
		repo:  repo,
	}
}

func (targ *OrderService) GetOrder(uid string) (model.Order, error) {
	order, err := targ.cache.Get(uid)
	if err == nil {
		log.Printf("Order %q found in cache", uid)
		return order, nil
	}

	if errors.Is(err, cache.ErrNotFound) {
		log.Printf("Order %q not in cache. Fetching from DB...", uid)

		orderFromDB, dbErr := targ.repo.GetByUID(uid)
		if dbErr != nil {
			return model.Order{}, dbErr
		}

		log.Printf("Order %q found in DB. Caching...", uid)
		targ.cache.Set(orderFromDB)

		return orderFromDB, nil
	}
	return model.Order{}, err
}
