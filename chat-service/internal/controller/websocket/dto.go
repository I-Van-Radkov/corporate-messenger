package websocket

import (
	"encoding/json"
	"time"
)

// WSMessage - общая структура WS сообщения
type WSMessage struct {
	Type      string          `json:"type"` // send_message, typing, read_receipt, etc.
	Payload   json.RawMessage `json:"payload"`
	Timestamp int64           `json:"timestamp"`
	ChatID    string          `json:"chat_id,omitempty"`
}

// WSSendMessagePayload - payload для отправки сообщения
type WSSendMessagePayload struct {
	Content string `json:"content"`
	Type    string `json:"type"` // text, file, image, system
	ReplyTo string `json:"reply_to,omitempty"`
}

// WSEditMessagePayload - payload для редактирования сообщения
type WSEditMessagePayload struct {
	MessageID string `json:"message_id"`
	Content   string `json:"content"`
}

// WSDeleteMessagePayload - payload для удаления сообщения
type WSDeleteMessagePayload struct {
	MessageID string `json:"message_id"`
}

// WSReadReceiptPayload - payload для отметки прочитанным
type WSReadReceiptPayload struct {
	MessageID string `json:"message_id"`
	ChatID    string `json:"chat_id"`
}

// WSTypingPayload - payload для индикатора набора текста
type WSTypingPayload struct {
	ChatID   string `json:"chat_id"`
	IsTyping bool   `json:"is_typing"`
}

// WSStatusUpdatePayload - payload для обновления статуса
type WSStatusUpdatePayload struct {
	Status    string `json:"status"` // online, away, offline, dnd
	Timestamp int64  `json:"timestamp"`
}

// WSEventMessage - структура для рассылки событий
type WSEventMessage struct {
	Event     string      `json:"event"` // message_sent, message_edited, user_joined, etc.
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}
