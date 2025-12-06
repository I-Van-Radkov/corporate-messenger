package models

import (
	"time"

	"github.com/google/uuid"
)

type AccountRole string

const (
	RoleUser      AccountRole = "user"
	RoleSupport   AccountRole = "support"
	RoleAdmin     AccountRole = "admin"
	RoleModerator AccountRole = "moderator"
)

type Account struct {
	AccountID    uuid.UUID   `json:"account_id" db:"account_id"`
	UserID       uuid.UUID   `json:"user_id" db:"user_id"`
	Email        string      `json:"email" db:"email"`
	PasswordHash string      `json:"-" db:"password_hash"` // Не сериализуется в JSON
	Role         AccountRole `json:"role" db:"role"`
	IsActive     bool        `json:"is_active" db:"is_active"`
	LastLogin    *time.Time  `json:"last_login,omitempty" db:"last_login"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
}
