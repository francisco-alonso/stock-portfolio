package api

import (
	"encoding/json"
	"net/http"

	"github.com/francisco-alonso/stock-portfolio/portfolio-service/services"
)

// Handler represents the API controller.
type Handler struct {
    portfolioService services.PortfolioService
}

// NewHandler initializes a new handler.
func NewHandler(portfolioService services.PortfolioService) *Handler {
    return &Handler{portfolioService: portfolioService}
}

// GetPortfolio devuelve la cartera de un usuario.
func (h *Handler) GetPortfolio(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("user_id")
    if userID == "" {
        http.Error(w, "user_id is required", http.StatusBadRequest)
        return
    }

    positions, err := h.portfolioService.GetPositions(userID)
    if err != nil {
        http.Error(w, "Error retrieving portfolio", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(positions)
}

// AddPosition handles adding/updating positions.
func (h *Handler) AddPosition(w http.ResponseWriter, r *http.Request) {
    var req PositionDto
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    err := h.portfolioService.AddPosition(req.UserID, req.Asset, req.Quantity, req.Price)
    if err != nil {
        http.Error(w, "Failed to save position", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "Position saved successfully"})
}
