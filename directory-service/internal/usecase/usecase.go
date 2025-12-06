package usecase

import (
	"context"

	"github.com/I-Van-Radkov/corporate-messenger/directory-service/internal/models"
	"github.com/google/uuid"
)

type DirectoryRepo interface {
	CreateDepartment(ctx context.Context, dep *models.Department) (uuid.UUID, error)
	CreateUser(ctx context.Context, user *models.User) (uuid.UUID, error)
	GetDepartmentById(ctx context.Context, id uuid.UUID) (*models.Department, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error)
}

type DirectoryUsecase struct {
	Repository DirectoryRepo
}

func NewDirectoryUsecase(repository DirectoryRepo) *DirectoryUsecase {
	return &DirectoryUsecase{
		Repository: repository,
	}
}
