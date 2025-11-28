package models

import (
	"time"

	"github.com/google/uuid"
)

type Department struct {
	DepartmentID uuid.UUID  `json:"department_id" db:"department_id"`
	Name         string     `json:"name" db:"name"`
	ParentID     *uuid.UUID `json:"parent_id,omitempty" db:"parent_id"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`

	Parent   *Department   `json:"parent,omitempty" db:"-"`
	Children []*Department `json:"children,omitempty" db:"-"`
	Users    []*User       `json:"users,omitempty" db:"-"`
}
