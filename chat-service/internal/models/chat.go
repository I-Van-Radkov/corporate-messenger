package models

import (
	"time"

	"github.com/google/uuid"
)

type ChatType string
type MessageType string
type MemberRole string

const (
	ChatPrivate    ChatType = "private"
	ChatGroup      ChatType = "group"
	ChatDepartment ChatType = "department"

	MsgText   MessageType = "text"
	MsgFile   MessageType = "file"
	MsgImage  MessageType = "image"
	MsgSystem MessageType = "system"

	RoleOwner  MemberRole = "owner"
	RoleAdmin  MemberRole = "admin"
	RoleMember MemberRole = "member"
)

type Chat struct {
	ChatID    uuid.UUID `json:"chat_id" db:"chat_id"`
	Type      ChatType  `json:"type" db:"type"`
	Name      *string   `json:"name,omitempty" db:"name"`
	CreatedBy uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ChatMember struct {
	ChatID   uuid.UUID  `json:"chat_id" db:"chat_id"`
	UserID   uuid.UUID  `json:"user_id" db:"user_id"`
	Role     MemberRole `json:"role" db:"role"`
	JoinedAt time.Time  `json:"joined_at" db:"joined_at"`
}

type Message struct {
	MessageID  uuid.UUID   `json:"message_id" db:"message_id"`
	ChatID     uuid.UUID   `json:"chat_id" db:"chat_id"`
	SenderID   uuid.UUID   `json:"sender_id" db:"sender_id"`
	Content    string      `json:"content" db:"content"`
	Type       MessageType `json:"type" db:"type"`
	ReplyTo    *uuid.UUID  `json:"reply_to,omitempty" db:"reply_to"`
	IsEdited   bool        `json:"is_edited" db:"is_edited"`
	IsDeleted  bool        `json:"is_deleted" db:"is_deleted"`
	SentAt     time.Time   `json:"sent_at" db:"sent_at"`
	ReadCount  int         `json:"read_count,omitempty" db:"read_count"`
	IsReadByMe bool        `json:"is_read_by_me,omitempty" db:"-"`
}
