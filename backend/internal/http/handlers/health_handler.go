package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Checker interface {
	Ping(ctx context.Context) error
}

type HealthHandler struct {
	db Checker
}

func NewHealthHandler(db Checker) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) RegisterRoutes(r chi.Router) {
	r.Get("/live", h.Liveness)
	r.Get("/ready", h.Readiness)
}

func (h *HealthHandler) Liveness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "UP",
	})
}

func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	// Verificar conexión a la base de datos
	if err := h.db.Ping(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status": "DOWN",
			"error":  "database connection failed",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "UP",
	})
}
