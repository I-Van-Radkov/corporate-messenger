package adapter

import (
	"context"

	"github.com/I-Van-Radkov/corporate-messenger/identity-service/internal/models"
	"github.com/Masterminds/squirrel"
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

func (r *AuthRepo) Create(ctx context.Context, account *models.Account) error
func (r *AuthRepo) FindByID(ctx context.Context, accountID string) (*models.Account, error)
func (r *AuthRepo) FindByEmail(ctx context.Context, email string) (*models.Account, error)
func (r *AuthRepo) FindByUserID(ctx context.Context, userID string) (*models.Account, error)
