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
	"github.com/Jiangultimo/notification-bridge/internal/handler"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	port := ":8000"
	if v := os.Getenv("PORT"); v != "" {
		port = ":" + v
	}

	mux := http.NewServeMux()

	// docs endpoint
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("/healthz", handler.Healthz)

	srv := &http.Server{
		Addr:              port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on %s", port)
	log.Fatal(srv.ListenAndServe())

}
