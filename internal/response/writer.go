package response

import (
	"encoding/json"
	"net/http"
)

const contentTypeJSON = "application/json; charset=utf-8"

type requestIDGetter interface {
	RequestID() string
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func OK[T any](w http.ResponseWriter, requestId string, data T) {
	resp := Envelope[T]{
		Code:      0,
		Success:   true,
		Message:   "ok",
		Data:      &data,
		RequestID: requestId,
	}
	WriteJSON(w, http.StatusOK, resp)
}

func Fail(w http.ResponseWriter, status int, requestID string, code int, message string) {
	resp := Envelope[any]{
		Code:      code,
		Success:   false,
		Message:   message,
		Data:      nil,
		RequestID: requestID,
	}
	WriteJSON(w, status, resp)
}
