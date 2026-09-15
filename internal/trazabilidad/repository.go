package trazabilidad

import (
	"context"

	"github.com/google/uuid"
)

type Appender interface {
	Append(ctx context.Context, e Evento) error
}

type Repository interface {
	Appender
	ListByCarga(ctx context.Context, orgID, cargaID uuid.UUID) ([]Evento, error)
}
