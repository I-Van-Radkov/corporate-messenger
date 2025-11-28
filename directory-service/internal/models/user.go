package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	Email        string     `json:"email" db:"email"`
	FirstName    string     `json:"first_name" db:"first_name"`
	LastName     string     `json:"last_name" db:"last_name"`
	Position     *string    `json:"position,omitempty" db:"position"`
	DepartmentID *uuid.UUID `json:"department_id,omitempty" db:"department_id"`
	AvatarURL    *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`

	Department *Department `json:"department,omitempty" db:"-"`
}
