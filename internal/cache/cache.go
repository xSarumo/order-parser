package cache

import (
	"container/list"
	"database/sql"
	"errors"
	"log"
	"sync"
	"test-task/internal/config"
	"test-task/internal/model"
	"test-task/internal/repository"
)

var CACHE_LIMIT int = config.CacheLimit()

var ErrNotFound = errors.New("cache: item not found")

type LRU_Cache struct {
	mtx        sync.RWMutex
	linkedList *list.List
	storage    map[string]*list.Element
}

type entry struct {
	key   string
	value model.Order
}

func NewCache() *LRU_Cache {
	return &LRU_Cache{
		linkedList: list.New(),
		storage:    make(map[string]*list.Element),
	}
}

func (targ *LRU_Cache) Set(value model.Order) {
	targ.mtx.Lock()
	defer targ.mtx.Unlock()
	key := value.OrderUID
	if elem, hit := targ.storage[key]; hit {
		targ.linkedList.MoveToFront(elem)
		targ.storage[key].Value.(*entry).value = value
		return
	}
	newElem := entry{key, value}
	newListElement := targ.linkedList.PushFront(&newElem)
	targ.storage[key] = newListElement

	if targ.linkedList.Len() > CACHE_LIMIT {
		if oldest := targ.linkedList.Back(); oldest != nil {
			targ.linkedList.Remove(oldest)
			delete(targ.storage, oldest.Value.(*entry).key)
		}
	}
}

func (targ *LRU_Cache) Get(key string) (model.Order, error) {
	targ.mtx.Lock()
	defer targ.mtx.Unlock()

	if elem, hit := targ.storage[key]; hit {
		targ.linkedList.MoveToFront(elem)
		return elem.Value.(*entry).value, nil
	}
	return model.Order{}, ErrNotFound
}

func (targ *LRU_Cache) LoadFromDB(orders []model.Order) {
	for _, order := range orders {
		targ.Set(order)
	}
}

func InitCache(db *sql.DB) *LRU_Cache {
	cache := NewCache()
	limit := config.CacheLimit()
	orders, err := repository.NewOrderRepository(db).GetLastNOrders(limit)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	cache.LoadFromDB(orders)
	return cache
}
