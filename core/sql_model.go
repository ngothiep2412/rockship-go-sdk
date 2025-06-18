package core

import (
	"time"

	"github.com/google/uuid"
)

type SQLModel struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func NewSQLModel() SQLModel {
	now := time.Now().UTC()

	return SQLModel{
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}
