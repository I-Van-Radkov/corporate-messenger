package dto

import "github.com/google/uuid"

type UserResponse struct {
	UserID       uuid.UUID  `json:"user_id"`
	Email        string     `json:"email"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	IsActive     bool       `json:"is_active"`
	DepartmentID *uuid.UUID `json:"department,omitempty"`
}
