package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Jiangultimo/notification-bridge/internal/types"
)

// BarkClient 是 Bark API 客户端
type BarkClient struct {
	BaseURL string
	Key     string
	Icon    string
	client  *http.Client
}

// NewBarkClient 创建一个新的 Bark 客户端
func NewBarkClient(baseURL, key, icon string) *BarkClient {
	return &BarkClient{
		BaseURL: baseURL,
		Key:     key,
		Icon:    icon,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendNotification 发送通知到 Bark
func (c *BarkClient) SendNotification(req *types.BarkRequest) error {
	// 构建完整的 URL
	url := fmt.Sprintf("%s/%s", c.BaseURL, c.Key)

	// 如果请求中没有指定 icon,使用客户端默认的
	if req.Icon == "" && c.Icon != "" {
		req.Icon = c.Icon
	}

	// 序列化请求体
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// 创建 HTTP 请求
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json; charset=utf-8")

	// 发送请求
	log.Printf("[Bark] Sending notification to %s", url)
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var barkResp types.BarkResponse
	if err := json.NewDecoder(resp.Body).Decode(&barkResp); err != nil {
		log.Printf("[Bark] Failed to decode response: %v", err)
		// 即使解析失败，如果 HTTP 状态码是成功的，也认为请求成功
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}
		return fmt.Errorf("failed to decode response: %w", err)
	}

	// 检查响应状态
	if resp.StatusCode >= 200 && resp.StatusCode < 300 && barkResp.Code == 200 {
		log.Printf("[Bark] Notification sent successfully: %s", barkResp.Message)
		return nil
	}

	return fmt.Errorf("bark API returned error: code=%d, message=%s", barkResp.Code, barkResp.Message)
}
