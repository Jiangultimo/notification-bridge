// @title Notification Bridge API Server
// @version 1.0.0
// @description Main entry point for the Notification Bridge API server.
// @BasePath /
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/Jiangultimo/notification-bridge/docs"
	"github.com/Jiangultimo/notification-bridge/internal/client"
	"github.com/Jiangultimo/notification-bridge/internal/handler"
	"github.com/Jiangultimo/notification-bridge/internal/handler/webhook"
	"github.com/Jiangultimo/notification-bridge/internal/middleware"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

// initBarkClient 初始化 Bark 客户端
func initBarkClient() *client.BarkClient {
	barkAPIURL := os.Getenv("BARK_API_URL")
	barkKey := os.Getenv("BARK_KEY")
	barkIcon := os.Getenv("BARK_ICON")

	if barkAPIURL != "" && barkKey != "" {
		log.Printf("Initializing Bark client with URL: %s", barkAPIURL)
		return client.NewBarkClient(barkAPIURL, barkKey, barkIcon)
	}
	log.Printf("Bark configuration missing (BARK_API_URL or BARK_KEY), notifications will not be sent")
	return nil
}

func main() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// 初始化 Bark 客户端
	barkClient := initBarkClient()

	port := ":8000"
	if v := os.Getenv("PORT"); v != "" {
		port = ":" + v
	}

	mux := http.NewServeMux()

	// docs endpoint
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("/healthz", handler.Healthz)
	mux.HandleFunc("/api/webhook/netalertx", webhook.NewNetAlertXHandler(barkClient))

	handler := middleware.WithRequestID(mux)
	srv := &http.Server{
		Addr:              port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on %s", port)
	log.Fatal(srv.ListenAndServe())

}
