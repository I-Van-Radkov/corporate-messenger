package dto

import (
	"time"

	"github.com/I-Van-Radkov/corporate-messenger/chat-service/internal/models"
)

// CreateChatRequest - запрос на создание чата
type CreateChatRequest struct {
	Type      string   `json:"type" validate:"required,oneof=private group department"`
	Name      *string  `json:"name" validate:"required_if=Type group,required_if=Type department"` // обязательное для group/department
	MemberIDs []string `json:"member_ids" validate:"required,min=1"`
}

// AddMembersRequest - запрос на добавление участников
type AddMembersRequest struct {
	UserIDs []string `json:"user_ids" validate:"required,min=1"`
}

// ChangeRoleRequest - запрос на смену роли
type ChangeRoleRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
	Role   string `json:"role" validate:"required,oneof=owner admin member"`
}

// FileInfoDTO - информация о файле
type FileInfoDTO struct {
	ID   string `json:"id"`
	URL  string `json:"url"`
	Name string `json:"name"`
	Size int64  `json:"size"`
	Type string `json:"type"`
}

// ChatPreview - превью чата для списка
type ChatPreview struct {
	ID           string          `json:"id"`
	Type         models.ChatType `json:"type"`
	Name         *string         `json:"name,omitempty"`
	AvatarURL    *string         `json:"avatar_url,omitempty"`
	MemberCount  int             `json:"member_count"`
	LastMessage  *MessageDTO     `json:"last_message,omitempty"`
	UnreadCount  int             `json:"unread_count"`
	LastActivity time.Time       `json:"last_activity"`
	IsPinned     bool            `json:"is_pinned"`
	IsMuted      bool            `json:"is_muted"`
}

// GetMessagesRequest - запрос на получение сообщений
type GetMessagesRequest struct {
	Limit  int     `form:"limit" validate:"max=100"`
	Before *string `form:"before" validate:"omitempty,uuid"` // cursor
	After  *string `form:"after" validate:"omitempty,uuid"`
}

// MessagesPageDTO - страница сообщений
type MessagesPageDTO struct {
	Messages      []MessageDTO `json:"messages"`
	HasMoreBefore bool         `json:"has_more_before"`
	HasMoreAfter  bool         `json:"has_more_after"`
}

// ChatDetailDTO - детальная информация о чате
type ChatDetailDTO struct {
	ID        string          `json:"id"`
	Type      models.ChatType `json:"type"`
	Name      *string         `json:"name,omitempty"`
	CreatedBy string          `json:"created_by"`
	CreatedAt time.Time       `json:"created_at"`
	Members   []ChatMemberDTO `json:"members"`
	PinnedIDs []string        `json:"pinned_ids,omitempty"`
}
