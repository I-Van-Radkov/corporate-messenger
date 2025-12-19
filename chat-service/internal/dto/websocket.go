package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type RequestType string

type EventType string

type ErrorType string

const (
	// Входящие типы (от клиента)
	TypeSendMessage RequestType = "send_message"

	// Исходящие типы (к клиенту)
	EventMessageSent EventType = "message.sent"
	EventError       EventType = "error"

	// Типы ошибок
	ErrSendMsg          ErrorType = "send_message_error"
	ErrInvalidMsgFormat ErrorType = "ivalid_message_format"
	ErrInvalidMsgType   ErrorType = "invalid_message_type"
	ErrInvalidPayload   ErrorType = "invalid_payload"
	ErrDataIsEmpty      ErrorType = "data_is_empty"
	ErrAccessDenied     ErrorType = "access_denied"
	ErrSaveFailed       ErrorType = "save_failed"
	ErrInternalError    ErrorType = "internal_error"
)

// IncomingMessage - входящее сообщение от клиента
type IncomingMessage struct {
	Type    RequestType     `json:"type"`    // Тип операции
	Payload json.RawMessage `json:"payload"` // Сырые данные\
	ChatID  uuid.UUID       `json:"chat_id"`
}

type SendMessageIncPayload struct {
	Content string     `json:"content"`
	Type    string     `json:"type"` // "text", "image", "file"
	ReplyTo *uuid.UUID `json:"reply_to,omitempty"`
}

// OutgoingMessage - исходящее сообщение к клиенту
type OutgoingMessage struct {
	Type    EventType       `json:"type"`    // Тип события
	Payload json.RawMessage `json:"payload"` // Сериализованные данные
	Meta    MessageMeta     `json:"meta"`    // Метаданные
}

// MessageSentPayload - событие отправки сообщения
type MessageSentOutPayload struct {
	Message MessageDTO `json:"message"`
}

type MessageMeta struct {
	Timestamp time.Time `json:"timestamp"`
	EventID   uuid.UUID `json:"event_id,omitempty"`
	ChatID    uuid.UUID `json:"chat_id,omitempty"`
}

// ErrorPayload - сообщение об ошибке
type ErrorPayload struct {
	Code    ErrorType `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
}

// MessageDTO - DTO сообщения
type MessageDTO struct {
	ID        uuid.UUID  `json:"id"`
	ChatID    uuid.UUID  `json:"chat_id"`
	SenderID  uuid.UUID  `json:"sender_id"`
	Content   string     `json:"content"`
	Type      string     `json:"type"`
	ReplyTo   *uuid.UUID `json:"reply_to,omitempty"`
	SentAt    time.Time  `json:"sent_at"`
	EditedAt  *time.Time `json:"edited_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
