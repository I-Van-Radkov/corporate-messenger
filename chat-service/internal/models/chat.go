package models

import (
	"time"

	"github.com/google/uuid"
)

type ChatType string

const (
	ChatTypePrivate    ChatType = "private"
	ChatTypeGroup      ChatType = "group"
	ChatTypeDepartment ChatType = "department"
)

type MemberRole string

const (
	MemberRoleOwner  MemberRole = "owner"
	MemberRoleAdmin  MemberRole = "admin"
	MemberRoleMember MemberRole = "member"
)

type MessageType string

const (
	MessageTypeText   MessageType = "text"
	MessageTypeFile   MessageType = "file"
	MessageTypeImage  MessageType = "image"
	MessageTypeSystem MessageType = "system"
)

// Chat - структура чата
type Chat struct {
	ID        uuid.UUID `db:"chat_id"`
	Type      ChatType  `db:"type"`
	Name      *string   `db:"name"` // NULL для private-чатов
	CreatedBy uuid.UUID `db:"created_by"`
	CreatedAt time.Time `db:"created_at"`
}

// ChatMember - участник чата
type ChatMember struct {
	ChatID   uuid.UUID  `db:"chat_id"`
	UserID   uuid.UUID  `db:"user_id"`
	Role     MemberRole `db:"role"`
	JoinedAt time.Time  `db:"joined_at"`
}

// Message - сообщение
type Message struct {
	ID        uuid.UUID   `db:"message_id"`
	ChatID    uuid.UUID   `db:"chat_id"`
	SenderID  uuid.UUID   `db:"sender_id"`
	Content   string      `db:"content"`
	Type      MessageType `db:"type"`
	ReplyTo   *uuid.UUID  `db:"reply_to"` // NULL если не ответ
	IsEdited  bool        `db:"is_edited"`
	IsDeleted bool        `db:"is_deleted"`
	SentAt    time.Time   `db:"sent_at"`
}

// MessageRead - прочитанные сообщения
type MessageRead struct {
	MessageID uuid.UUID `db:"message_id"`
	UserID    uuid.UUID `db:"user_id"`
	ReadAt    time.Time `db:"read_at"`
}

// UnreadCount - счётчик непрочитанных
type UnreadCount struct {
	UserID            uuid.UUID  `db:"user_id"`
	ChatID            uuid.UUID  `db:"chat_id"`
	UnreadCount       int        `db:"unread_count"`
	LastReadMessageID *uuid.UUID `db:"last_read_message_id"`
}

// PinnedMessage - закреплённые сообщения
type PinnedMessage struct {
	ChatID    uuid.UUID `db:"chat_id"`
	MessageID uuid.UUID `db:"message_id"`
	PinnedAt  time.Time `db:"pinned_at"`
	PinnedBy  uuid.UUID `db:"pinned_by"`
}
