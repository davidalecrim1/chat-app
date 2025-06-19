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
	connectedUsers: make(map[types.UserId]*types.UserConnection),
	mu:             sync.RWMutex{},
}

type connectionStore struct {
	connectedUsers map[types.UserId]*types.UserConnection
	mu             sync.RWMutex
}

func (store *connectionStore) add(conn *types.UserConnection, id types.UserId) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.connectedUsers[id] = conn
}

func (store *connectionStore) delete(id types.UserId) {
	store.mu.Lock()
	defer store.mu.Unlock()
	delete(store.connectedUsers, id)
}

func (store *connectionStore) get(id types.UserId) *types.UserConnection {
	return store.connectedUsers[id]
}

func (store *connectionStore) broadcastMessage(msg []byte) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	for _, connectedUser := range store.connectedUsers {
		err := connectedUser.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			slog.Error("failed to send message", "RemoteAddr", connectedUser.RemoteAddr(), "error", err)
		}
	}
}

func (store *connectionStore) broadcastDisconnectedUser(uConn *types.UserConnection) {
	msg := types.DisconnectedUserMessage{
		Name: uConn.User.Name,
	}

	rawMsg, err := json.Marshal(msg)
	if err != nil {
		slog.Error("failed to encode message", "error", err)
	}

	env := types.MessageEnvelope{
		Type:    "disconnectedUser",
		Payload: rawMsg,
	}

	rawEnv, err := json.Marshal(env)
	if err != nil {
		slog.Error("failed to encode message", "error", err)
	}

	store.broadcastMessage(rawEnv)
}

func (store *connectionStore) broadcastConnectedUser(uConn *types.UserConnection) {
	msg := types.ConnectedUserMessage{
		Name: uConn.User.Name,
	}

	rawMsg, err := json.Marshal(msg)
	if err != nil {
		slog.Error("failed to encode message", "error", err)
	}

	env := types.MessageEnvelope{
		Type:    "connectedUser",
		Payload: rawMsg,
	}

	rawEnv, err := json.Marshal(env)
	if err != nil {
		slog.Error("failed to encode message", "error", err)
	}

	store.broadcastMessage(rawEnv)
}

func JoinChat(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Info("Failed to upgrade the connection to websocket", "error", err)
		return
	}

	id := r.URL.Query().Get("id")
	name := r.URL.Query().Get("name")

	if id == "" {
		slog.Warn("Missing userId in WebSocket URL")
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(
			websocket.ClosePolicyViolation, "userId is required",
		))
		conn.Close()
		return
	}

	userId := types.UserId(id)
	userConnected := types.NewUserConnection(conn, types.User{Id: userId, Name: name})
	connStore.add(userConnected, userId)
	connStore.broadcastConnectedUser(userConnected)
	readMessagesFromConnection(conn, userId)
}

func readMessagesFromConnection(conn *websocket.Conn, id types.UserId) {
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				slog.Info("WebSocket closed normally", slog.String("error", err.Error()))
				break
			}
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
				slog.Info("WebSocket closed unexpectedly", slog.String("error", err.Error()))
				break
			}
			if closeErr, ok := err.(*websocket.CloseError); ok && closeErr.Code == 1005 {
				slog.Info("WebSocket closed without status (1005)", slog.String("error", err.Error()))
				break
			}

			slog.Info("WebSocket read error", slog.String("error", err.Error()))
			break
		}

		handleMessages(raw)
	}

	deletedUser := connStore.get(id)
	connStore.delete(id)
	connStore.broadcastDisconnectedUser(deletedUser)
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
			slog.Error("invalid message provided, ignoring it.", "error", err, "raw", raw)
			return
		}

		connStore.broadcastMessage(raw)
	}
}

func EnableBroadcastActiveConnections() {
	for {
		connStore.mu.RLock()
		activeConns := len(connStore.connectedUsers)
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

		time.Sleep(time.Second * 10)
	}
}
