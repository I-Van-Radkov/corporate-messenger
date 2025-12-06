package adapter

import (
	"context"

	"github.com/I-Van-Radkov/corporate-messenger/chat-service/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRepo struct {
	db      *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewChatRepo(db *pgxpool.Pool) *ChatRepo {
	return &ChatRepo{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *ChatRepo) GetChatMembers(ctx context.Context, chatID string) ([]*models.ChatMember, error)

func (r *ChatRepo) AddMembers(ctx context.Context, chatID string, userIDs []string, adderID string) error
func (r *ChatRepo) RemoveMember(ctx context.Context, chatID, userID, removerID string) error
func (r *ChatRepo) ChangeMemberRole(ctx context.Context, chatID, userID string, newRole models.MemberRole, changerID string) error

func (r *ChatRepo) GetUserChats(ctx context.Context, userID string, limit, offset int) ([]*models.Chat, error)

func (r *ChatRepo) SendMessage(ctx context.Context, msg *models.Message) (uuid.UUID, error)
func (r *ChatRepo) EditMessage(ctx context.Context, messageID, userID uuid.UUID, content string) error
func (r *ChatRepo) DeleteMessage(ctx context.Context, messageID, userID uuid.UUID) error
func (r *ChatRepo) MarkAsRead(ctx context.Context, userID, chatID, messageID uuid.UUID) error
func (r *ChatRepo) IsUserInChat(ctx context.Context, userID, chatID uuid.UUID) (bool, error)
