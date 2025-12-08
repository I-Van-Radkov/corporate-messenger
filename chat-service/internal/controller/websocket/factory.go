package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/I-Van-Radkov/corporate-messenger/chat-service/internal/dto"
	"github.com/google/uuid"
)

type MessageFactory struct {
}

func NewMessageFactory() *MessageFactory {
	return &MessageFactory{}
}

func (f *MessageFactory) NewOutgoingMessage(msg *dto.MessageDTO) *dto.OutgoingMessage {
	payload := &dto.MessageSentOutPayload{
		Message: *msg,
	}

	return f.createOutgoingMessage(dto.EventMessageSent, payload, msg.ChatID)
}

func (f *MessageFactory) NewError(code dto.ErrorType, message, details string) *dto.OutgoingMessage {
	payload := dto.ErrorPayload{
		Code:    code,
		Message: message,
		Details: details,
	}

	return f.createOutgoingMessage(dto.EventError, payload, uuid.Nil)
}

func (f *MessageFactory) createOutgoingMessage(eventtype dto.EventType, payload interface{}, chatID uuid.UUID) *dto.OutgoingMessage {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal payload: %v", err)

		return &dto.OutgoingMessage{
			Type:    dto.EventError,
			Payload: []byte(`{"code":"internal_error","message":"Failed to create event"}`),
			Meta: dto.MessageMeta{
				Timestamp: time.Now(),
				EventID:   uuid.New(),
				ChatID:    chatID,
			},
		}
	}

	return &dto.OutgoingMessage{
		Type:    eventtype,
		Payload: payloadBytes,
		Meta: dto.MessageMeta{
			Timestamp: time.Now(),
			EventID:   uuid.New(),
			ChatID:    chatID,
		},
	}
}
