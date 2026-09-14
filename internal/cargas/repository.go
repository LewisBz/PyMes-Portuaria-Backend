package cargas

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, c *Carga) error
	List(ctx context.Context, orgID uuid.UUID, clienteID *uuid.UUID, estado *Estado) ([]Carga, error)
	GetByID(ctx context.Context, orgID, id uuid.UUID) (*Carga, error)
	Update(ctx context.Context, c *Carga) error
}
