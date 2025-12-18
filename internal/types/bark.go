package types

// BarkRequest 是发送到 Bark API 的请求结构
type BarkRequest struct {
	Title    string `json:"title"`
	Body     string `json:"body,omitempty"`
	Markdown string `json:"markdown,omitempty"`
	Badge    int    `json:"badge,omitempty"`
	Icon     string `json:"icon,omitempty"`
	Group    string `json:"group,omitempty"`
	Sound    string `json:"sound,omitempty"`
}

// BarkResponse 是 Bark API 的响应结构
type BarkResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}
