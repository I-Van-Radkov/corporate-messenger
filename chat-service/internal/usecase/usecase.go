package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/I-Van-Radkov/corporate-messenger/chat-service/internal/dto"
	"github.com/I-Van-Radkov/corporate-messenger/chat-service/internal/errors"
	"github.com/I-Van-Radkov/corporate-messenger/chat-service/internal/models"
	"github.com/google/uuid"
)

type ChatRepo interface {
	CreateChat(ctx context.Context, chat *models.Chat, memberIDs []uuid.UUID) error
	GetChat(ctx context.Context, chatId uuid.UUID) (*models.Chat, error)
	RemoveChat(ctx context.Context, chatId uuid.UUID) error
	GetUserChats(ctx context.Context, userID string) ([]*models.Chat, error)

	GetLastMessageInChat(ctx context.Context, chatId uuid.UUID) (*models.Message, error)
	GetChatMessages(ctx context.Context, chatID, userID uuid.UUID, limit int, before *uuid.UUID) ([]*models.Message, error)

	GetChatMembers(ctx context.Context, chatID uuid.UUID) ([]*models.ChatMember, error)
	AddMembers(ctx context.Context, chatID uuid.UUID, userIDs []uuid.UUID, adderID uuid.UUID) error
	RemoveMember(ctx context.Context, chatID, userID, removerID uuid.UUID) error
	ChangeMemberRole(ctx context.Context, chatID, userID uuid.UUID, newRole models.MemberRole, changerID uuid.UUID) error

	//GetUserChats(ctx context.Context, userID string, limit, offset int) ([]*models.Chat, error)

	SendMessage(ctx context.Context, msg *models.Message) (uuid.UUID, error)
	IsUserInChat(ctx context.Context, userID, chatID uuid.UUID) (bool, error)

	// EditMessage(ctx context.Context, messageID, userID uuid.UUID, content string) error
	// DeleteMessage(ctx context.Context, messageID, userID uuid.UUID) error
	// MarkAsRead(ctx context.Context, userID, chatID, messageID uuid.UUID) error
}

type OfflineMessageStorage interface {
	SendMessage(userId uuid.UUID, data []byte)
	GetMessages(userId uuid.UUID) ([][]byte, error)
}

type ChatUsecase struct {
	chatRepo   ChatRepo
	msgStorage OfflineMessageStorage
}

func NewChatUsecase(chatRepo ChatRepo, msgStorage OfflineMessageStorage) *ChatUsecase {
	return &ChatUsecase{
		chatRepo: chatRepo,
	}
}

func (u *ChatUsecase) SendMsgToStorage(ctx context.Context, userId uuid.UUID, msg []byte) {
	u.msgStorage.SendMessage(userId, msg)
}

func (u *ChatUsecase) GetMessagesFromStorage(ctx context.Context, userId uuid.UUID) [][]byte {
	messages, err := u.msgStorage.GetMessages(userId)
	if err != nil {
		return nil
	}

	return messages
}

func (u *ChatUsecase) SendMessageToDb(ctx context.Context, msg *dto.MessageDTO) (string, error) {
	message := &models.Message{
		MessageID: uuid.New(),
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

func (u *ChatUsecase) IsUserInChat(ctx context.Context, userID, chatID uuid.UUID) (bool, error) {
	return u.chatRepo.IsUserInChat(ctx, userID, chatID)
}

func (u *ChatUsecase) CreateChat(ctx context.Context, req *dto.CreateChatRequest, creatorID uuid.UUID) (*dto.CreateChatResponse, error) {
	if req.Type == dto.ChatPrivate && len(req.MemberIDs) != 1 {
		return nil, fmt.Errorf("private chat must have exactly one member")
	}
	if (req.Type == dto.ChatGroup || req.Type == dto.ChatDepartment) && req.Name == nil {
		return nil, fmt.Errorf("name required for group/department")
	}

	var memberIDs []uuid.UUID
	memberIDs = append(memberIDs, creatorID) // создатель всегда участник
	memberIDs = append(memberIDs, req.MemberIDs...)

	chat := &models.Chat{
		ChatID:    uuid.New(),
		Type:      models.ChatType(req.Type),
		Name:      req.Name,
		CreatedBy: creatorID,
		CreatedAt: time.Now(),
	}

	err := u.chatRepo.CreateChat(ctx, chat, memberIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat: %w", err)
	}

	resp := &dto.CreateChatResponse{
		ChatID: chat.ChatID,
	}

	return resp, nil
}
func (u *ChatUsecase) RemoveChat(ctx context.Context, chatID uuid.UUID) error {
	_, err := u.chatRepo.GetChat(ctx, chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrChatNotFound
		}

		return fmt.Errorf("failed to get chat: %w", err)
	}

	err = u.chatRepo.RemoveChat(ctx, chatID)
	if err != nil {
		return fmt.Errorf("failed to remove chat: %w", err)
	}

	return nil
}

func (u *ChatUsecase) GetUserChats(ctx context.Context, userID uuid.UUID) (*dto.GetUserChatsResponse, error) {
	chats, err := u.chatRepo.GetUserChats(ctx, uuid.NewString())
	if err != nil {
		return nil, fmt.Errorf("failed to get chats: %w", err)
	}

	total := len(chats)
	chatsPreviews := make([]dto.ChatPreview, total)
	for i := 0; i < total; i++ {
		chat := chats[i]

		lastMsg, err := u.chatRepo.GetLastMessageInChat(ctx, chat.ChatID)
		if err != nil {
			return nil, fmt.Errorf("failed to get last message in chat: %w", err)
		}

		lastMsgPreview := &dto.MessagePreview{
			MessageID: lastMsg.MessageID,
			Content:   lastMsg.Content,
			Type:      dto.MessageType(lastMsg.Type),
			SenderID:  lastMsg.SenderID,
			SentAt:    lastMsg.SentAt,
		}

		chatsPreviews[i] = dto.ChatPreview{
			ChatID:          chat.ChatID,
			Type:            dto.ChatType(chat.Type),
			Name:            chat.Name,
			LastMessage:     lastMsgPreview,
			LastMessageTime: &lastMsg.SentAt,
		}
	}

	resp := &dto.GetUserChatsResponse{
		Total: total,
		Chats: chatsPreviews,
	}

	return resp, nil
}
func (u *ChatUsecase) GetChatMessages(ctx context.Context, chatID, userID uuid.UUID, limit int, before *uuid.UUID) (*dto.GetMessagesResponse, error) {
	_, err := u.chatRepo.GetChat(ctx, chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrChatNotFound
		}

		return nil, fmt.Errorf("failed to get chat: %w", err)
	}

	ok, err := u.chatRepo.IsUserInChat(ctx, userID, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user in chat: %w", err)
	}
	if !ok {
		return nil, errors.ErrUserNotInChat
	}

	messages, err := u.chatRepo.GetChatMessages(ctx, chatID, userID, limit, before)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages from chat^ %w", err)
	}

	total := len(messages)
	msgResp := make([]dto.MessageDTO, total)
	for i := 0; i < total; i++ {
		message := messages[i]
		msgResp[i] = dto.MessageDTO{
			ID:       message.MessageID,
			ChatID:   message.ChatID,
			SenderID: message.SenderID,
			Content:  message.Content,
			Type:     string(message.Type),
			ReplyTo:  message.ReplyTo,
			SentAt:   message.SentAt,
		}
	}

	resp := &dto.GetMessagesResponse{
		Messages: msgResp,
		Total:    total,
	}

	return resp, nil
}

func (u *ChatUsecase) GetChatMembers(ctx context.Context, chatID uuid.UUID) (*dto.ChatMembers, error) {
	members, err := u.chatRepo.GetChatMembers(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get members from db: %w", err)
	}

	total := len(members)
	membersDTO := make([]dto.ChatMemberDTO, total)
	for i := 0; i < total; i++ {
		member := members[i]
		membersDTO[i] = dto.ChatMemberDTO{
			UserID:   member.UserID,
			ChatID:   member.ChatID,
			Role:     string(member.Role),
			JoinedAt: member.JoinedAt,
		}
	}

	return &dto.ChatMembers{
		Members: membersDTO,
		Total:   total,
	}, nil
}

func (u *ChatUsecase) AddMembers(ctx context.Context, chatID uuid.UUID, userIDs []uuid.UUID, adderID uuid.UUID) error {
	_, err := u.chatRepo.GetChat(ctx, chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrChatNotFound
		}

		return fmt.Errorf("failed to get chat: %w", err)
	}

	err = u.chatRepo.AddMembers(ctx, chatID, userIDs, adderID)
	if err != nil {
		return fmt.Errorf("failed to save members to db: %w", err)
	}

	return nil
}

func (u *ChatUsecase) RemoveMember(ctx context.Context, chatID, userID, removerID uuid.UUID) error {
	_, err := u.chatRepo.GetChat(ctx, chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrChatNotFound
		}

		return fmt.Errorf("failed to get chat: %w", err)
	}

	ok, err := u.chatRepo.IsUserInChat(ctx, userID, chatID)
	if err != nil {
		return fmt.Errorf("failed to get user in chat: %w", err)
	}
	if !ok {
		return errors.ErrUserNotInChat
	}

	err = u.chatRepo.RemoveMember(ctx, chatID, userID, removerID)
	if err != nil {
		return fmt.Errorf("failed to remove member from db: %w", err)
	}

	return nil
}

func (u *ChatUsecase) ChangeMemberRole(ctx context.Context, chatID, userID uuid.UUID, role dto.MemberRole, changerID uuid.UUID) error {
	_, err := u.chatRepo.GetChat(ctx, chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrChatNotFound
		}

		return fmt.Errorf("failed to get chat: %w", err)
	}

	ok, err := u.chatRepo.IsUserInChat(ctx, userID, chatID)
	if err != nil {
		return fmt.Errorf("failed to get user in chat: %w", err)
	}
	if !ok {
		return errors.ErrUserNotInChat
	}

	err = u.chatRepo.ChangeMemberRole(ctx, changerID, userID, models.MemberRole(role), changerID)
	if err != nil {
		return fmt.Errorf("failed to change member role from db: %w", err)
	}

	return nil
}
