package directory

import "errors"

type UserResponse struct {
	UserID     string `json:"user_id"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	IsActive   bool   `json:"is_active"`
	Department string `json:"department,omitempty"`
}

var (
	ErrUserNotFound = errors.New("user not found")
)
