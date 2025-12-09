package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/I-Van-Radkov/corporate-messenger/directory-service/internal/dto"
	"github.com/I-Van-Radkov/corporate-messenger/directory-service/internal/errors"
	"github.com/I-Van-Radkov/corporate-messenger/directory-service/internal/models"
	"github.com/google/uuid"
)

type DirectoryRepo interface {
	CreateDepartment(ctx context.Context, dep *models.Department) (uuid.UUID, error)
	GetDepartmentByID(ctx context.Context, id uuid.UUID) (*models.Department, error)
	GetDepartments(ctx context.Context, limit, offset int) ([]*models.Department, int, error)
	DeleteDepartment(ctx context.Context, id uuid.UUID) error
	GetDepartmentMembers(ctx context.Context, depID uuid.UUID, limit, offset int) ([]*models.User, int, error)

	CreateUser(ctx context.Context, user *models.User) (uuid.UUID, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUsers(ctx context.Context, filter *dto.GetUsersRequest) ([]*models.User, int, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error

	// Для дерева отделов
	GetAllDepartments(ctx context.Context) ([]*models.Department, error)
}

type DirectoryUsecase struct {
	Repository DirectoryRepo
}

func NewDirectoryUsecase(repository DirectoryRepo) *DirectoryUsecase {
	return &DirectoryUsecase{
		Repository: repository,
	}
}

func (u *DirectoryUsecase) CreateDepartment(ctx context.Context, req *dto.CreateDepartmentRequest) (*dto.CreateDepartmentResponse, error) {
	now := time.Now()

	dep := &models.Department{
		DepartmentID: uuid.New(),
		Name:         req.Name,
		ParentID:     req.ParentID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if req.ParentID != nil {
		parent, err := u.Repository.GetDepartmentByID(ctx, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get parent department: %w", err)
		}
		if parent == nil {
			return nil, errors.ErrDepartmentNotFound
		}

		if parent.ParentID != nil && *parent.ParentID == dep.DepartmentID {
			return nil, errors.ErrCircularReference
		}
	}

	id, err := u.Repository.CreateDepartment(ctx, dep)
	if err != nil {
		return nil, fmt.Errorf("failed to save department: %w", err)
	}

	resp := &dto.CreateDepartmentResponse{
		DepartmentID: id,
	}

	return resp, nil
}

func (u *DirectoryUsecase) GetDepartments(ctx context.Context, req *dto.GetDepartmentsRequest) (*dto.GetDepartmentsResponse, error) {
	if req.Tree {
		allDeps, err := u.Repository.GetAllDepartments(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get all departments")
		}
		tree := buildDepartmentTree(allDeps)

		return &dto.GetDepartmentsResponse{
			Departments: tree,
			Total:       len(allDeps),
		}, nil
	}

	deps, total, err := u.Repository.GetDepartments(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all departments")
	}

	return &dto.GetDepartmentsResponse{
		Departments: deps,
		Total:       total,
	}, nil
}

func buildDepartmentTree(deps []*models.Department) []*models.Department {
	depMap := make(map[uuid.UUID]*models.Department)
	for _, dep := range deps {
		depMap[dep.DepartmentID] = dep
		dep.Children = []*models.Department{}
	}

	var roots []*models.Department
	for _, dep := range deps {
		if dep.ParentID == nil {
			roots = append(roots, dep)
		} else {
			parent, ok := depMap[*dep.ParentID]
			if ok {
				parent.Children = append(parent.Children, dep)
			}
		}
	}
	return roots
}

func (u *DirectoryUsecase) RemoveDepartment(ctx context.Context, id uuid.UUID) error {
	dep, err := u.Repository.GetDepartmentByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get department by id: %w", err)
	}
	if dep == nil {
		return errors.ErrDepartmentNotFound
	}

	return u.Repository.DeleteDepartment(ctx, id)
}

func (u *DirectoryUsecase) GetDepartmentMembers(ctx context.Context, req *dto.GetDepartmentMembersRequest) (*dto.GetDepartmentMembersResponse, error) {
	dep, err := u.Repository.GetDepartmentByID(ctx, req.DepartmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get department by id: %w", err)
	}
	if dep == nil {
		return nil, errors.ErrDepartmentNotFound
	}

	members, total, err := u.Repository.GetDepartmentMembers(ctx, req.DepartmentID, req.Limit, req.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get members of department: %w", err)
	}

	return &dto.GetDepartmentMembersResponse{
		Users: members,
		Total: total,
	}, nil
}

func (u *DirectoryUsecase) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
	now := time.Now()

	user := &models.User{
		UserID:       uuid.New(),
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Position:     req.Position,
		DepartmentID: req.DepartmentID,
		AvatarURL:    req.AvatarURL,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if req.DepartmentID != nil {
		dep, err := u.Repository.GetDepartmentByID(ctx, *req.DepartmentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get department by id: %w", err)
		}
		if dep == nil {
			return nil, errors.ErrDepartmentNotFound
		}
	}

	id, err := u.Repository.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return &dto.CreateUserResponse{
		UserID: id,
	}, nil
}

func (u *DirectoryUsecase) GetUsers(ctx context.Context, req *dto.GetUsersRequest) (*dto.GetUsersResponse, error) {
	if req.DepartmentID != nil {
		dep, err := u.Repository.GetDepartmentByID(ctx, *req.DepartmentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get department by id: %w", err)
		}
		if dep == nil {
			return nil, errors.ErrDepartmentNotFound
		}
	}

	users, total, err := u.Repository.GetUsers(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return &dto.GetUsersResponse{
		Users: users,
		Total: total,
	}, nil
}

func (u *DirectoryUsecase) GetUser(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	user, err := u.Repository.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}
	return &dto.UserResponse{
		UserID:       user.UserID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		IsActive:     user.IsActive,
		DepartmentID: user.DepartmentID,
	}, nil
}

func (u *DirectoryUsecase) RemoveUser(ctx context.Context, id uuid.UUID) error {
	user, err := u.Repository.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrUserNotFound
	}
	return u.Repository.DeleteUser(ctx, id)
}
