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
	User User
	conn *websocket.Conn
	mu   sync.RWMutex
}

func NewUserConnection(conn *websocket.Conn, u User) *UserConnection {
	return &UserConnection{
		User: u,
		conn: conn,
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
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ConnectedUserMessage struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
