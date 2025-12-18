package webhook

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Jiangultimo/notification-bridge/internal/client"
	"github.com/Jiangultimo/notification-bridge/internal/middleware"
	"github.com/Jiangultimo/notification-bridge/internal/response"
	"github.com/Jiangultimo/notification-bridge/internal/types"
)

// NetAlertXResponse 是接口的响应数据
type NetAlertXResponse struct {
	Received bool   `json:"received"`
	Message  string `json:"message"`
}

// NetAlertXEnvelope 用于 Swagger 文档
type NetAlertXEnvelope struct {
	Code      int               `json:"code"`
	Success   bool              `json:"success"`
	Message   string            `json:"message"`
	Data      NetAlertXResponse `json:"data"`
	RequestID string            `json:"request_id,omitempty"`
}

// NewNetAlertXHandler 创建一个新的 NetAlertX webhook handler
func NewNetAlertXHandler(barkClient *client.BarkClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		netAlertXHandler(w, r, barkClient)
	}
}

// netAlertXHandler 是实际的处理函数
// NetAlertX handles webhook requests from NetAlertX
// @Summary NetAlertX Webhook Endpoint
// @Description Receives webhook notifications from NetAlertX, logs and processes the data
// @Tags Webhook
// @Accept json
// @Produce json
// @Param payload body types.NetAlertXPayload true "Webhook payload"
// @Success 200 {object} NetAlertXEnvelope "Successfully received"
// @Failure 400 {object} response.ErrorEnvelope "Bad Request"
// @Failure 405 {object} response.ErrorEnvelope "Method Not Allowed"
// @Router /api/webhook/netalertx [post]
func netAlertXHandler(w http.ResponseWriter, r *http.Request, barkClient *client.BarkClient) {
	rid := middleware.GetRequestID(r)

	// 只接受 POST 方法
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		response.Fail(w, http.StatusMethodNotAllowed, rid, 1, "method not allowed")
		return
	}

	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[%s] Failed to read request body: %v", rid, err)
		response.Fail(w, http.StatusBadRequest, rid, 2, "failed to read request body")
		return
	}
	defer r.Body.Close()

	// 解析外层 JSON
	var payload types.NetAlertXPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("[%s] Failed to parse JSON: %v", rid, err)
		response.Fail(w, http.StatusBadRequest, rid, 3, "invalid JSON format")
		return
	}

	// 格式化打印日志
	log.Printf("[%s] ========== NetAlertX Webhook ==========", rid)
	log.Printf("[%s] Username: %s", rid, payload.Username)
	log.Printf("[%s] Text: %s", rid, payload.Text)

	// 解析 attachments
	if len(payload.Attachments) > 0 {
		attachment := payload.Attachments[0]
		log.Printf("[%s] Title: %s", rid, attachment.Title)

		// 解析内部 JSON 字符串
		var data types.NetAlertXData
		if err := json.Unmarshal([]byte(attachment.Text), &data); err != nil {
			log.Printf("[%s] Failed to parse attachment text: %v", rid, err)
		} else {
			printNetAlertXData(rid, &data)

			// 转发到 Bark
			if barkClient != nil {
				markdown := formatToMarkdown(&data)
				barkReq := &types.BarkRequest{
					Title:    "NetAlertX",
					Markdown: markdown,
					Badge:    1,
					Group:    "NetAlertX",
				}

				// 异步发送，不阻塞响应
				go func() {
					if err := barkClient.SendNotification(barkReq); err != nil {
						log.Printf("[%s] Failed to send Bark notification: %v", rid, err)
					}
				}()
			}
		}
	}

	log.Printf("[%s] =======================================", rid)

	// 返回成功响应
	respData := NetAlertXResponse{
		Received: true,
		Message:  "Webhook received and logged successfully",
	}
	response.OK(w, rid, respData)
}

// printNetAlertXData 格式化打印 NetAlertX 数据
func printNetAlertXData(rid string, data *types.NetAlertXData) {
	// 打印新设备
	if len(data.NewDevices) > 0 {
		log.Printf("[%s] --- New Devices (%d) ---", rid, len(data.NewDevices))
		for i, dev := range data.NewDevices {
			log.Printf("[%s]   [%d] MAC: %s, IP: %s, Name: %s, Time: %s",
				rid, i+1, dev.MAC, dev.IP, dev.DeviceName, dev.Datetime)
		}
	}

	// 打印重连设备
	if len(data.DownReconnected) > 0 {
		log.Printf("[%s] --- Down Reconnected (%d) ---", rid, len(data.DownReconnected))
		for i, dev := range data.DownReconnected {
			log.Printf("[%s]   [%d] Name: %s, MAC: %s, IP: %s, Vendor: %s",
				rid, i+1, dev.DevName, dev.MAC, dev.IP, dev.Vendor)
		}
	}

	// 打印离线设备
	if len(data.DownDevices) > 0 {
		log.Printf("[%s] --- Down Devices (%d) ---", rid, len(data.DownDevices))
		for i, dev := range data.DownDevices {
			log.Printf("[%s]   [%d] Name: %s, MAC: %s, IP: %s, Vendor: %s",
				rid, i+1, dev.DevName, dev.MAC, dev.IP, dev.Vendor)
		}
	}

	// 打印事件
	if len(data.Events) > 0 {
		log.Printf("[%s] --- Events (%d) ---", rid, len(data.Events))
		for i, evt := range data.Events {
			log.Printf("[%s]   [%d] MAC: %s, IP: %s, Type: %s, Time: %s",
				rid, i+1, evt.MAC, evt.IP, evt.EventType, evt.Datetime)
		}
	}

	// 打印插件信息
	if len(data.Plugins) > 0 {
		log.Printf("[%s] --- Plugins (%d) ---", rid, len(data.Plugins))
		for i, plugin := range data.Plugins {
			log.Printf("[%s]   [%d] Plugin: %s, MAC: %s, IP: %s, Status: %s",
				rid, i+1, plugin.Plugin, plugin.ObjectPrimaryID,
				plugin.WatchedValue1, plugin.Status)
		}
	}
}

// formatToMarkdown 将 NetAlertX 数据格式化为 Markdown
func formatToMarkdown(data *types.NetAlertXData) string {
	var builder strings.Builder

	// 新设备
	if len(data.NewDevices) > 0 {
		builder.WriteString(fmt.Sprintf("### 🆕 新设备 (%d)\n", len(data.NewDevices)))
		for _, dev := range data.NewDevices {
			builder.WriteString(fmt.Sprintf("- **设备**: %s\n", dev.DeviceName))
			builder.WriteString(fmt.Sprintf("  - **MAC**: %s\n", dev.MAC))
			builder.WriteString(fmt.Sprintf("  - **IP**: %s\n", dev.IP))
			builder.WriteString(fmt.Sprintf("  - **时间**: %s\n", dev.Datetime))
			builder.WriteString("\n")
		}
	}

	// 重连设备
	if len(data.DownReconnected) > 0 {
		builder.WriteString(fmt.Sprintf("### 🔁 重连设备 (%d)\n", len(data.DownReconnected)))
		for _, dev := range data.DownReconnected {
			builder.WriteString(fmt.Sprintf("- **设备**: %s\n", dev.DevName))
			builder.WriteString(fmt.Sprintf("  - **MAC**: %s\n", dev.MAC))
			builder.WriteString(fmt.Sprintf("  - **IP**: %s\n", dev.IP))
			builder.WriteString(fmt.Sprintf("  - **厂商**: %s\n", dev.Vendor))
			builder.WriteString(fmt.Sprintf("  - **时间**: %s\n", dev.DateTime))
			builder.WriteString("\n")
		}
	}

	// 离线设备
	if len(data.DownDevices) > 0 {
		builder.WriteString(fmt.Sprintf("### 🔴 离线设备 (%d)\n", len(data.DownDevices)))
		for _, dev := range data.DownDevices {
			builder.WriteString(fmt.Sprintf("- **设备**: %s\n", dev.DevName))
			builder.WriteString(fmt.Sprintf("  - **MAC**: %s\n", dev.MAC))
			builder.WriteString(fmt.Sprintf("  - **IP**: %s\n", dev.IP))
			builder.WriteString(fmt.Sprintf("  - **厂商**: %s\n", dev.Vendor))
			builder.WriteString(fmt.Sprintf("  - **时间**: %s\n", dev.DateTime))
			builder.WriteString("\n")
		}
	}

	// 事件
	if len(data.Events) > 0 {
		builder.WriteString(fmt.Sprintf("### 📋 事件 (%d)\n", len(data.Events)))
		for _, evt := range data.Events {
			builder.WriteString(fmt.Sprintf("- **MAC**: %s | **IP**: %s\n", evt.MAC, evt.IP))
			builder.WriteString(fmt.Sprintf("  - **类型**: %s\n", evt.EventType))
			builder.WriteString(fmt.Sprintf("  - **时间**: %s\n", evt.Datetime))
			builder.WriteString("\n")
		}
	}

	// 插件事件
	if len(data.Plugins) > 0 {
		builder.WriteString(fmt.Sprintf("### 🔌 插件事件 (%d)\n", len(data.Plugins)))
		for _, plugin := range data.Plugins {
			builder.WriteString(fmt.Sprintf("- **插件**: %s\n", plugin.Plugin))
			builder.WriteString(fmt.Sprintf("  - **MAC**: %s\n", plugin.ObjectPrimaryID))
			builder.WriteString(fmt.Sprintf("  - **IP**: %s\n", plugin.WatchedValue1))
			builder.WriteString(fmt.Sprintf("  - **状态**: %s\n", plugin.Status))
			builder.WriteString("\n")
		}
	}

	result := builder.String()
	if result == "" {
		return "暂无事件"
	}
	return result
}
