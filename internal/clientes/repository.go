package clientes

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, c *Cliente) error
	List(ctx context.Context, orgID uuid.UUID) ([]Cliente, error)
	GetByID(ctx context.Context, orgID, id uuid.UUID) (*Cliente, error)
}
