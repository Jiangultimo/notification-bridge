package handler

import (
	"net/http"
	"time"

	"github.com/Jiangultimo/notification-bridge/internal/middleware"
	"github.com/Jiangultimo/notification-bridge/internal/response"
)

type HealthResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

type HealthEnvelope struct {
	Code      int            `json:"code"`
	Success   bool           `json:"success"`
	Message   string         `json:"message"`
	Data      HealthResponse `json:"data"`
	RequestID string         `json:"request_id,omitempty"`
}

// Healthz handles the /healthz endpoint to report service health status.
// @Summary Health Check Endpoint
// @Description Returns the health status of the service.
// @Tags Health
// @Produce json
// @Success 200 {object} HealthEnvelope "OK"
// @Failure 405 {object} response.ErrorEnvelope "Method Not Allowed"
// @Router /healthz [get]
func Healthz(w http.ResponseWriter, r *http.Request) {
	rid := middleware.GetRequestID(r)

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		response.Fail(w, http.StatusMethodNotAllowed, rid, 1, "method not allowed")
		return
	}

	data := HealthResponse{
		Status: "ok",
		Time:   time.Now().UTC().Format(time.RFC3339),
	}
	response.OK(w, rid, data)
}
