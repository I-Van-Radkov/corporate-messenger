package dto

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
)

type CreateChatRequest struct {
	Type      ChatType    `json:"type" binding:"required,oneof=private group department"`
	Name      *string     `json:"name,omitempty"`
	MemberIDs []uuid.UUID `json:"member_ids" binding:"required,min=1,dive,required"`
}

type CreateChatResponse struct {
	ChatID uuid.UUID `json:"chat_id"`
}

type GetUserChatsResponse struct {
	Chats []ChatPreview `json:"chats"`
	Total int           `json:"total"`
}

type ChatPreview struct {
	ChatID uuid.UUID `json:"chat_id"`
	Type   ChatType  `json:"type"`
	Name   *string   `json:"name,omitempty"`
	//UnreadCount     int             `json:"unread_count"`
	LastMessage     *MessagePreview `json:"last_message,omitempty"`
	LastMessageTime *time.Time      `json:"last_message_time,omitempty"`
}

type MessagePreview struct {
	MessageID uuid.UUID   `json:"message_id"`
	Content   string      `json:"content"`
	Type      MessageType `json:"type"`
	SenderID  uuid.UUID   `json:"sender_id"`
	SentAt    time.Time   `json:"sent_at"`
}

type GetMessagesRequest struct {
	ChatID uuid.UUID  `uri:"chat_id" binding:"required"`
	Before *uuid.UUID `form:"before,omitempty"`
	Limit  int        `form:"limit,default=50"`
}

type GetMessagesResponse struct {
	Messages []MessageDTO `json:"messages"`
	Total    int          `json:"total"`
}

// type SendMessageRequest struct {
// 	Content  string            `json:"content" binding:"required"`
// 	Type     models.MessageType `json:"type" binding:"required,oneof=text file image system"`
// 	ReplyTo  *uuid.UUID        `json:"reply_to,omitempty"`
// }

// type SendMessageResponse struct {
// 	MessageID uuid.UUID `json:"message_id"`
// 	SentAt    time.Time `json:"sent_at"`
// }

type AddMembersRequest struct {
	UserIDs []uuid.UUID `json:"user_ids" binding:"required,min=1,dive,required"`
}

type ChangeRoleRequest struct {
	Role MemberRole `json:"role" binding:"required,oneof=admin member"`
}

type ChatMemberDTO struct {
	UserID   uuid.UUID `json:"user_id"`
	ChatID   uuid.UUID `json:"chat_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

type ChatMembers struct {
	Members []ChatMemberDTO `json:"members"`
	Total   int             `json:"total"`
}
