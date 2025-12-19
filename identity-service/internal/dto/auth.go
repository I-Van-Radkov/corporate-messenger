package dto

import "time"

type AccountRole string

const (
	RoleUser      AccountRole = "user"
	RoleSupport   AccountRole = "support"
	RoleAdmin     AccountRole = "admin"
	RoleModerator AccountRole = "moderator"
)

type CreateAccountRequest struct {
	UserID   string      `json:"user_id" validate:"required,uuid"`
	Email    string      `json:"email" validate:"required,email"`
	Password string      `json:"password" validate:"required,min=8"`
	Role     AccountRole `json:"role" validate:"required,oneof=user support admin moderator"`
}

type AccountResponse struct {
	AccountID string      `json:"account_id"`
	UserID    string      `json:"user_id"`
	Email     string      `json:"email"`
	Role      AccountRole `json:"role"`
	IsActive  bool        `json:"is_active"`
	LastLogin *time.Time  `json:"last_login,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type IntrospectRequest struct {
	Token string `json:"token"`
}

type IntrospectResponse struct {
	Active bool `json:"active"`
}

type UpdateAccountRequest struct {
	Email    *string     `json:"email,omitempty" validate:"omitempty,email"`
	Role     AccountRole `json:"role,omitempty" validate:"omitempty,oneof=user support admin moderator"`
	IsActive *bool       `json:"is_active,omitempty"`
}
