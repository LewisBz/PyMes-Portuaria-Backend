package organizacion

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Organizacion, error)
	IsActiva(ctx context.Context, id uuid.UUID) (bool, error)
}
