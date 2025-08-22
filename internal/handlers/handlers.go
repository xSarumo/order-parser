package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"test-task/internal/service"

	"github.com/go-chi/chi/v5"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderUID := chi.URLParam(r, "order_uid")
	if orderUID == "" {
		http.Error(w, "Order UID is required", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrder(orderUID)
	if err != nil {
		log.Printf("Failed to get order %s: %v", orderUID, err)
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		log.Printf("Failed to encode order %s to JSON: %v", orderUID, err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}
