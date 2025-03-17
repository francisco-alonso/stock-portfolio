package api

import (
	"encoding/json"
	"net/http"

	"github.com/francisco-alonso/trade-engine-service/internal/application"
	"github.com/francisco-alonso/trade-engine-service/internal/domain"
)

type Handler struct {
	service *application.TradeService
}

func NewHandler(service *application.TradeService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/create-order", h.CreateOrder)
	return mux
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order domain.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid order format", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.service.CreateOrder(ctx, order); err != nil {
		http.Error(w, "Failed to publish order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Order published successfully"))
}