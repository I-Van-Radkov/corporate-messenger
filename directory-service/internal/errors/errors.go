package errors

import "errors"

var (
	// Department errors
	ErrDepartmentNotFound    = errors.New("department not found")
	ErrDepartmentHasUsers    = errors.New("department has active users")
	ErrDepartmentHasChildren = errors.New("department has child departments")
	ErrCircularReference     = errors.New("circular reference detected")
	ErrSelfParentReference   = errors.New("cannot set department as its own parent")

	// User errors
	ErrUserNotFound    = errors.New("user not found")
	ErrUserEmailExists = errors.New("user with this email already exists")

	// Common errors
	ErrInvalidUUID      = errors.New("invalid uuid format")
	ErrValidationFailed = errors.New("validation failed")
	ErrForbidden        = errors.New("access forbidden")
)
