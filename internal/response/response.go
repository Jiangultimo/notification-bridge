package response

type Envelope[T any] struct {
	Code      int    `json:"code"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Data      *T     `json:"data,omitempty"` // 使用指针和omitempty：可以在错误时不返回data字段
	RequestID string `json:"request_id,omitempty"`
}

// ErrorEnvelope is used for swagger documentation of error responses
type ErrorEnvelope struct {
	Code      int    `json:"code"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}
