package main

import (
	"log"
	"log/slog"
	"net/http"

	"chat-app/internal/handler"
)

func main() {
	http.HandleFunc("GET /ws/connect", handler.JoinChat)

	slog.Info("Starting server on port 8201...")
	go handler.EnableBroadcastActiveConnections()

	err := http.ListenAndServe(":8201", nil)
	if err != nil {
		log.Fatalf("failed to start server...")
	}
}
