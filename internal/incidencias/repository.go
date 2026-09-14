package incidencias

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, i *Incidencia) error
	ListByCarga(ctx context.Context, orgID, cargaID uuid.UUID) ([]Incidencia, error)
	GetByID(ctx context.Context, orgID, id uuid.UUID) (*Incidencia, error)
	Update(ctx context.Context, i *Incidencia) error
}
