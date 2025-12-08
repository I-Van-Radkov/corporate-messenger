package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type TypeMsg string

const (
	typingType  TypeMsg = "typing"
	sendMsgType TypeMsg = "send_message"
)

type WSMessage struct {
	Type      TypeMsg         `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp int64           `json:"timestamp"`
	ChatID    uuid.UUID       `json:"chat_id,omitempty"`
}
