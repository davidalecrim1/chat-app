package types

import "encoding/json"

type UserId string

type User struct {
	Id   string
	Name string
}

type MessageEnvelope struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type ChatPayload struct {
	UserId string `json:"userId"`
	Text   string `json:"text"`
}

type ActiveUsersMessage struct {
	Total int `json:"total"`
}
