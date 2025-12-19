package adapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/I-Van-Radkov/corporate-messenger/chat-service/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (r *ChatRepo) CreateChat(ctx context.Context, chat *models.Chat, memberIDs []uuid.UUID) error {
	if len(memberIDs) == 0 {
		return errors.New("at least one member is required")
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // откатится, если не закоммитим

	// 1. Создаём чат — только то, что не дефолтится в БД
	chatSQL, chatArgs, _ := r.builder.Insert("chats").
		Columns("type", "name", "created_by").
		Values(chat.Type, chat.Name, chat.CreatedBy).
		Suffix("RETURNING chat_id, created_at").
		ToSql()

	err = tx.QueryRow(ctx, chatSQL, chatArgs...).Scan(&chat.ChatID, &chat.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert chat: %w", err)
	}

	// 2. Добавляем участников
	// Создатель чата (chat.CreatedBy) всегда становится owner'ом
	insertQB := r.builder.Insert("chat_members").
		Columns("chat_id", "user_id", "role", "joined_at")

	for _, userID := range memberIDs {
		role := models.RoleMember
		if userID == chat.CreatedBy {
			role = models.RoleOwner
		}
		insertQB = insertQB.Values(chat.ChatID, userID, role, time.Now())
	}

	membersSQL, membersArgs, _ := insertQB.ToSql()

	if _, err := tx.Exec(ctx, membersSQL, membersArgs...); err != nil {
		return fmt.Errorf("insert chat members: %w", err)
	}

	// 3. Всё ок — коммитим
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *ChatRepo) GetChat(ctx context.Context, chatID uuid.UUID) (*models.Chat, error) {
	query := r.builder.Select("chat_id", "type", "name", "created_by", "created_at").
		From("chats").
		Where(squirrel.Eq{"chat_id": chatID})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var chat models.Chat
	var name sql.NullString
	err = r.db.QueryRow(ctx, sqlStr, args...).Scan(
		&chat.ChatID,
		&chat.Type,
		&name,
		&chat.CreatedBy,
		&chat.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	if name.Valid {
		chat.Name = &name.String
	}

	return &chat, nil
}

func (r *ChatRepo) RemoveChat(ctx context.Context, chatID uuid.UUID) error {
	// Каскадное удаление через FOREIGN KEY должно удалить всё связанное
	// Но для надёжности можно удалять в правильном порядке
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Удаляем из chat_members
	_, err = tx.Exec(ctx, "DELETE FROM chat_members WHERE chat_id = $1", chatID)
	if err != nil {
		return err
	}

	// Удаляем из chats
	_, err = tx.Exec(ctx, "DELETE FROM chats WHERE chat_id = $1", chatID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *ChatRepo) GetUserChats(ctx context.Context, userIDStr string) ([]*models.Chat, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	query := r.builder.Select("c.chat_id", "c.type", "c.name", "c.created_by", "c.created_at").
		From("chats c").
		Join("chat_members cm ON c.chat_id = cm.chat_id").
		Where(squirrel.Eq{"cm.user_id": userID}).
		OrderBy("c.created_at DESC")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []*models.Chat
	for rows.Next() {
		var chat models.Chat
		var name sql.NullString
		err := rows.Scan(
			&chat.ChatID,
			&chat.Type,
			&name,
			&chat.CreatedBy,
			&chat.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if name.Valid {
			chat.Name = &name.String
		}
		chats = append(chats, &chat)
	}

	return chats, nil
}

func (r *ChatRepo) GetLastMessageInChat(ctx context.Context, chatID uuid.UUID) (*models.Message, error) {
	query := r.builder.Select("message_id", "chat_id", "sender_id", "content", "type", "reply_to", "is_edited", "is_deleted", "sent_at").
		From("messages").
		Where(squirrel.Eq{"chat_id": chatID, "is_deleted": false}).
		OrderBy("sent_at DESC").
		Limit(1)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var msg models.Message
	var replyTo sql.NullString
	err = r.db.QueryRow(ctx, sqlStr, args...).Scan(
		&msg.MessageID,
		&msg.ChatID,
		&msg.SenderID,
		&msg.Content,
		&msg.Type,
		&replyTo,
		&msg.IsEdited,
		&msg.IsDeleted,
		&msg.SentAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // нет сообщений — нормально
		}
		return nil, err
	}

	if replyTo.Valid {
		parsed, _ := uuid.Parse(replyTo.String)
		msg.ReplyTo = &parsed
	}

	return &msg, nil
}

func (r *ChatRepo) GetChatMessages(ctx context.Context, chatID, userID uuid.UUID, limit int, before *uuid.UUID) ([]*models.Message, error) {
	if limit <= 0 {
		limit = 50
	}

	query := r.builder.Select("m.message_id", "m.chat_id", "m.sender_id", "m.content", "m.type", "m.reply_to", "m.is_edited", "m.is_deleted", "m.sent_at").
		From("messages m").
		Where(squirrel.Eq{"m.chat_id": chatID, "m.is_deleted": false})

	if before != nil {
		query = query.Where(squirrel.Lt{"m.message_id": *before})
	}

	query = query.OrderBy("m.sent_at DESC").Limit(uint64(limit))

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		var msg models.Message
		var replyTo sql.NullString
		err := rows.Scan(
			&msg.MessageID,
			&msg.ChatID,
			&msg.SenderID,
			&msg.Content,
			&msg.Type,
			&replyTo,
			&msg.IsEdited,
			&msg.IsDeleted,
			&msg.SentAt,
		)
		if err != nil {
			return nil, err
		}
		if replyTo.Valid {
			parsed, _ := uuid.Parse(replyTo.String)
			msg.ReplyTo = &parsed
		}
		messages = append(messages, &msg)
	}

	// Возвращаем в хронологическом порядке (сначала старые)
	for i := len(messages)/2 - 1; i >= 0; i-- {
		opp := len(messages) - 1 - i
		messages[i], messages[opp] = messages[opp], messages[i]
	}

	return messages, nil
}

func (r *ChatRepo) GetChatMembers(ctx context.Context, chatID uuid.UUID) ([]*models.ChatMember, error) {
	query := r.builder.Select("chat_id", "user_id", "role", "joined_at").
		From("chat_members").
		Where(squirrel.Eq{"chat_id": chatID})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*models.ChatMember
	for rows.Next() {
		var member models.ChatMember
		err := rows.Scan(&member.ChatID, &member.UserID, &member.Role, &member.JoinedAt)
		if err != nil {
			return nil, err
		}
		members = append(members, &member)
	}

	return members, nil
}

func (r *ChatRepo) AddMembers(ctx context.Context, chatID uuid.UUID, userIDs []uuid.UUID, adderID uuid.UUID) error {
	if len(userIDs) == 0 {
		return nil
	}

	insert := r.builder.Insert("chat_members").
		Columns("chat_id", "user_id", "role", "joined_at")

	for _, userID := range userIDs {
		insert = insert.Values(chatID, userID, models.RoleMember, time.Now())
	}

	sqlStr, args, err := insert.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sqlStr, args...)
	return err
}

func (r *ChatRepo) RemoveMember(ctx context.Context, chatID, userID, removerID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		"DELETE FROM chat_members WHERE chat_id = $1 AND user_id = $2",
		chatID, userID)
	return err
}

func (r *ChatRepo) ChangeMemberRole(ctx context.Context, chatID, userID uuid.UUID, newRole models.MemberRole, changerID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		"UPDATE chat_members SET role = $1 WHERE chat_id = $2 AND user_id = $3",
		newRole, chatID, userID)
	return err
}

func (r *ChatRepo) SendMessage(ctx context.Context, msg *models.Message) (uuid.UUID, error) {
	insert := r.builder.Insert("messages").
		Columns("message_id", "chat_id", "sender_id", "content", "type", "reply_to", "is_edited", "is_deleted").
		Values(msg.MessageID, msg.ChatID, msg.SenderID, msg.Content, msg.Type, msg.ReplyTo, msg.IsEdited, msg.IsDeleted)

	sqlStr, args, err := insert.ToSql()
	if err != nil {
		return uuid.Nil, err
	}

	_, err = r.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return uuid.Nil, err
	}

	return msg.MessageID, nil
}

func (r *ChatRepo) IsUserInChat(ctx context.Context, userID, chatID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM chat_members WHERE chat_id = $1 AND user_id = $2)",
		chatID, userID).Scan(&exists)
	return exists, err
}
