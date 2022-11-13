package entities

import (
	"time"

	"github.com/google/uuid"
)

type Pet struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e *Pet) GenerateID() error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	e.ID = id.String()
	return nil
}

func (e *Pet) DeepCopy() *Pet {
	return &Pet{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
