package types

import (
	"encoding/json"
	"net"
	"sync"

	"github.com/gorilla/websocket"
)

type UserId string

type User struct {
	Id   UserId `json:"id"`
	Name string `json:"name"`
}

type UserConnection struct {
	conn *websocket.Conn
	User User
	mu   sync.RWMutex
}

func NewUserConnection(conn *websocket.Conn, u User) *UserConnection {
	return &UserConnection{
		conn: conn,
		User: u,
		mu:   sync.RWMutex{},
	}
}

func (u *UserConnection) WriteMessage(messageType int, data []byte) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.conn.WriteMessage(messageType, data)
}

func (u *UserConnection) RemoteAddr() net.Addr {
	return u.conn.RemoteAddr()
}

type MessageEnvelope struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type ChatPayload struct {
	User User   `json:"user"`
	Text string `json:"text"`
}

type ActiveUsersMessage struct {
	Total int `json:"total"`
}

type DisconnectedUserMessage struct {
	Name string `json:"name"`
}

type ConnectedUserMessage struct {
	Name string `json:"name"`
}
