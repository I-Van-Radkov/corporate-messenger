package dto

import (
	"github.com/I-Van-Radkov/corporate-messenger/directory-service/internal/models"
	"github.com/google/uuid"
)

// CreateDepartmentRequest - запрос на создание отдела
type CreateDepartmentRequest struct {
	Name     string     `json:"name" binding:"required"`
	ParentID *uuid.UUID `json:"parent_id,omitempty"`
}

// CreateDepartmentResponse - ответ на создание отдела
type CreateDepartmentResponse struct {
	DepartmentID uuid.UUID `json:"department_id"`
}

// GetDepartmentsRequest - запрос на список отделов (с пагинацией и опцией для дерева)
type GetDepartmentsRequest struct {
	Limit  int  `form:"limit,default=20"`
	Offset int  `form:"offset,default=0"`
	Tree   bool `form:"tree,default=false"` // Флаг для возврата иерархической структуры
}

// GetDepartmentsResponse - ответ на список отделов
type GetDepartmentsResponse struct {
	Departments []*models.Department `json:"departments"`
	Total       int                  `json:"total"`
}

// RemoveDepartmentRequest - запрос на удаление отдела (по ID в пути)
type RemoveDepartmentRequest struct {
	DepartmentID uuid.UUID `uri:"department_id" binding:"required"`
}

// GetDepartmentMembersRequest - запрос на членов отдела (с пагинацией)
type GetDepartmentMembersRequest struct {
	DepartmentID uuid.UUID `uri:"department_id" binding:"required"`
	Limit        int       `form:"limit,default=20"`
	Offset       int       `form:"offset,default=0"`
}

// GetDepartmentMembersResponse - ответ на членов отдела
type GetDepartmentMembersResponse struct {
	Users []*models.User `json:"users"`
	Total int            `json:"total"`
}

// CreateUserRequest - запрос на создание пользователя
type CreateUserRequest struct {
	Email        string     `json:"email" binding:"required,email"`
	FirstName    string     `json:"first_name" binding:"required"`
	LastName     string     `json:"last_name" binding:"required"`
	Position     *string    `json:"position,omitempty"`
	DepartmentID *uuid.UUID `json:"department_id,omitempty"`
	AvatarURL    *string    `json:"avatar_url,omitempty"`
}

// CreateUserResponse - ответ на создание пользователя
type CreateUserResponse struct {
	UserID uuid.UUID `json:"user_id"`
}

// GetUsersRequest - запрос на список пользователей (с пагинацией и фильтрами)
type GetUsersRequest struct {
	Limit        int        `form:"limit,default=20"`
	Offset       int        `form:"offset,default=0"`
	DepartmentID *uuid.UUID `form:"department_id,omitempty"`
	IsActive     *bool      `form:"is_active,omitempty"`
}

// GetUsersResponse - ответ на список пользователей
type GetUsersResponse struct {
	Users []*models.User `json:"users"`
	Total int            `json:"total"`
}

// GetUserRequest - запрос на получение пользователя (по ID в пути)
type GetUserRequest struct {
	UserID uuid.UUID `uri:"user_id" binding:"required"`
}

// RemoveUserRequest - запрос на удаление пользователя (по ID в пути)
type RemoveUserRequest struct {
	UserID uuid.UUID `uri:"user_id" binding:"required"`
}
