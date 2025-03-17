package api

import (
	"net/http"

	"github.com/francisco-alonso/go-template/internal/services"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.URL.Query().Get("username")
	email := r.URL.Query().Get("email")

	if username == "" || email == "" {
		http.Error(w, "Missing username or email", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateUser(username, email); err != nil {	
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}
