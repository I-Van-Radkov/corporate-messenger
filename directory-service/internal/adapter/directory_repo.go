package adapter

import (
	"context"

	"github.com/I-Van-Radkov/corporate-messenger/directory-service/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DirectoryRepo struct {
	db      *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewDirectoryRepo(db *pgxpool.Pool) *DirectoryRepo {
	return &DirectoryRepo{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *DirectoryRepo) CreateDepartment(ctx context.Context, dep *models.Department) (uuid.UUID, error)

func (r *DirectoryRepo) CreateUser(ctx context.Context, user *models.User) (uuid.UUID, error)

func (r *DirectoryRepo) GetDepartmentById(ctx context.Context, id uuid.UUID) (*models.Department, error)

func (r *DirectoryRepo) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error)

//TODO: Реализовать далее все CRUDL методы
