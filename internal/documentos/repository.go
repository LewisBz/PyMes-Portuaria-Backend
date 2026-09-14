package documentos

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, d *Documento) error
	ListByCarga(ctx context.Context, orgID, cargaID uuid.UUID, onlyVisible bool) ([]Documento, error)
	GetByID(ctx context.Context, orgID, id uuid.UUID) (*Documento, error)
}
