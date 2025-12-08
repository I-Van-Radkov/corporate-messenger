package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/I-Van-Radkov/corporate-messenger/chat-service/internal/dto"
	"github.com/I-Van-Radkov/corporate-messenger/chat-service/internal/models"
	"github.com/google/uuid"
)

type ChatRepo interface {
	GetChatMembers(ctx context.Context, chatID string) ([]*models.ChatMember, error)

	AddMembers(ctx context.Context, chatID string, userIDs []string, adderID string) error
	RemoveMember(ctx context.Context, chatID, userID, removerID string) error
	ChangeMemberRole(ctx context.Context, chatID, userID string, newRole models.MemberRole, changerID string) error

	GetUserChats(ctx context.Context, userID string, limit, offset int) ([]*models.Chat, error)

	SendMessage(ctx context.Context, msg *models.Message) (uuid.UUID, error)
	EditMessage(ctx context.Context, messageID, userID uuid.UUID, content string) error
	DeleteMessage(ctx context.Context, messageID, userID uuid.UUID) error
	MarkAsRead(ctx context.Context, userID, chatID, messageID uuid.UUID) error
	IsUserInChat(ctx context.Context, userID, chatID uuid.UUID) (bool, error)
}

type ChatUsecase struct {
	chatRepo ChatRepo
}

func NewChatUsecase(chatRepo ChatRepo) *ChatUsecase {
	return &ChatUsecase{
		chatRepo: chatRepo,
	}
}

func (u *ChatUsecase) SendMessage(ctx context.Context, msg *dto.MessageDTO) (string, error) {
	message := &models.Message{
		ID:        uuid.New(),
		ChatID:    msg.ChatID,
		SenderID:  msg.SenderID,
		Content:   msg.Content,
		Type:      models.MessageType(msg.Type),
		ReplyTo:   msg.ReplyTo,
		IsEdited:  false,
		IsDeleted: false,
		SentAt:    time.Now(),
	}

	id, err := u.chatRepo.SendMessage(ctx, message)
	if err != nil {
		return "", fmt.Errorf("failed to save message to db: %w", err)
	}

	return id.String(), nil
}

func (u *ChatUsecase) GetChatMembers(ctx context.Context, chatID string) ([]*dto.ChatMemberDTO, error) {
	members, err := u.chatRepo.GetChatMembers(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get members from db: %w", err)
	}

	var membersDTO []*dto.ChatMemberDTO
	for _, member := range members {
		membersDTO = append(membersDTO, &dto.ChatMemberDTO{
			UserID:   member.ChatID,
			Role:     string(member.Role),
			JoinedAt: member.JoinedAt,
		})
	}

	return membersDTO, nil
}

func (u *ChatUsecase) CreateChat(ctx context.Context, input *dto.CreateChatRequest) (string, error)

func (u *ChatUsecase) AddMembers(ctx context.Context, chatID string, userIDs []string, adderID string) error
func (u *ChatUsecase) RemoveMember(ctx context.Context, chatID, userID, removerID string) error
func (u *ChatUsecase) ChangeMemberRole(ctx context.Context, chatID, userID string, newRole string, changerID string) error

func (u *ChatUsecase) GetUserChats(ctx context.Context, userID string, limit, offset int) ([]*dto.ChatPreview, error)

func (u *ChatUsecase) EditMessage(ctx context.Context, messageID, userID, content string) error
func (u *ChatUsecase) DeleteMessage(ctx context.Context, messageID, userID string) error
func (u *ChatUsecase) MarkAsRead(ctx context.Context, userID, chatID, messageID string) error
func (u *ChatUsecase) IsUserInChat(ctx context.Context, userID, chatID string) (bool, error)
