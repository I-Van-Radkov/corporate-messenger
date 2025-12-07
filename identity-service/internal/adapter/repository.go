package adapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/I-Van-Radkov/corporate-messenger/identity-service/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	db      *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewAuthRepo(db *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *AuthRepo) Create(ctx context.Context, account *models.Account) error {
	query, args, err := r.builder.
		Insert("identity-db").
		Columns("account_id", "user_id", "email", "password_hash", "role", "is_active", "last_login", "created_at").
		Values(account.AccountID, account.UserID, account.Email, account.PasswordHash, account.Role, account.IsActive, account.LastLogin, account.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert subscription: %w", err)
	}

	return nil
}
func (r *AuthRepo) FindByID(ctx context.Context, accountID uuid.UUID) (*models.Account, error) {
	query, args, err := r.builder.
		Select("account_id", "user_id", "email", "password_hash", "role", "is_active", "last_login", "created_at").
		From("identity-db").
		Where(squirrel.Eq{"account_id": accountID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var account models.Account
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&account.AccountID, &account.UserID,
		&account.Email, &account.PasswordHash, &account.Role,
		&account.IsActive, &account.LastLogin, &account.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan account: %w", err)
	}

	return &account, nil
}

func (r *AuthRepo) FindByEmail(ctx context.Context, email string) (*models.Account, error) {
	query, args, err := r.builder.
		Select("account_id", "user_id", "email", "password_hash", "role", "is_active", "last_login", "created_at").
		From("identity-db").
		Where(squirrel.Eq{"email": email}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var account models.Account
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&account.AccountID, &account.UserID,
		&account.Email, &account.PasswordHash, &account.Role,
		&account.IsActive, &account.LastLogin, &account.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan account: %w", err)
	}

	return &account, nil
}

func (r *AuthRepo) FindByUserID(ctx context.Context, userID uuid.UUID) (*models.Account, error) {
	query, args, err := r.builder.
		Select("account_id", "user_id", "email", "password_hash", "role", "is_active", "last_login", "created_at").
		From("identity-db").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var account models.Account
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&account.AccountID, &account.UserID,
		&account.Email, &account.PasswordHash, &account.Role,
		&account.IsActive, &account.LastLogin, &account.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan account: %w", err)
	}

	return &account, nil
}
