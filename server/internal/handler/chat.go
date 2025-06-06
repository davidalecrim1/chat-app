package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"chat-app/types"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // TODO: Update this for production
	},
}

var connStore = &connectionStore{
	conns: make(map[string]*websocket.Conn),
	mu:    sync.RWMutex{},
}

type connectionStore struct {
	conns map[string]*websocket.Conn
	mu    sync.RWMutex
}

func (store *connectionStore) add(conn *websocket.Conn, id string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.conns[id] = conn
}

func (store *connectionStore) delete(id string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	delete(store.conns, id)
}

func (store *connectionStore) broadcastMessage(msg []byte) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	writeMu := sync.Mutex{}
	for _, conn := range store.conns {
		writeMu.Lock()
		err := conn.WriteMessage(websocket.TextMessage, msg)
		writeMu.Unlock()

		if err != nil {
			slog.Error("failed to send message", "RemoteAddr", conn.RemoteAddr(), "error", err)
		}
	}
}

func JoinChat(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Info("Failed to upgrade the connection to websocket", "error", err)
		return
	}

	id := r.URL.Query().Get("id")
	// name := r.URL.Query().Get("name")

	if id == "" {
		slog.Warn("Missing userId in WebSocket URL")
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(
			websocket.ClosePolicyViolation, "userId is required",
		))
		conn.Close()
		return
	}

	connStore.add(conn, id)
	readMessagesFromConnection(conn, id)
}

func readMessagesFromConnection(conn *websocket.Conn, id string) {
	for {
		messageType, raw, err := conn.ReadMessage()
		if err != nil && messageType != websocket.TextMessage {
			connStore.delete(id)
			slog.Error("received an invalid type of message", "error", err)
			return
		}

		handleMessages(raw)
	}
}

func handleMessages(raw []byte) {
	var envelope types.MessageEnvelope
	if err := json.Unmarshal(raw, &envelope); err != nil {
		slog.Error("failed to parse the message, the format is invalid.", "error", err, "raw", raw)
	}

	switch envelope.Type {
	case "chat":
		var p types.ChatPayload
		err := json.Unmarshal(envelope.Payload, &p)
		if err != nil {
			slog.Error("invalid message provided, ignoring it.", "error", err, "rawMsg", raw)
			return
		}

		connStore.broadcastMessage(raw)
	}
}

func EnableBroadcastActiveConnections() {
	for {
		connStore.mu.RLock()
		activeConns := len(connStore.conns)
		connStore.mu.RUnlock()

		if activeConns > 0 {
			msg := types.ActiveUsersMessage{
				Total: activeConns,
			}

			rawMsg, err := json.Marshal(msg)
			if err != nil {
				slog.Error("failed to create the message to broadcast", "error", err)
				return
			}

			envelope := types.MessageEnvelope{
				Type:    "activeUsers",
				Payload: rawMsg,
			}

			rawEnv, err := json.Marshal(envelope)
			if err != nil {
				slog.Error("failed to create message to broadcast amount of connections", "error", err)
				return
			}

			connStore.broadcastMessage(rawEnv)
		}

		time.Sleep(time.Second * 2)
	}
}
