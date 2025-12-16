package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

type HealthResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

// Healthz handles the /healthz endpoint to report service health status.
// @Summary Health Check Endpoint
// @Description Returns the health status of the service.
// @Tags Health
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 405 {string} string "Method Not Allowed"
// @Router /healthz [get]
func Healthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(HealthResponse{
		Status: "ok",
		Time:   time.Now().UTC().Format(time.RFC3339),
	})
}
