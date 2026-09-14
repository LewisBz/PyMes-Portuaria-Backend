package trazabilidad

import (
	"context"

	"github.com/google/uuid"
	"github.com/pymes-portuaria/backend/pkg/httperr"
)

type CargaGuard interface {
	EnsureVisible(ctx context.Context, orgID uuid.UUID, rol string, clienteScope *uuid.UUID, cargaID uuid.UUID) error
}

type Service struct {
	Repo  Repository
	Carga CargaGuard
}

func (s *Service) List(ctx context.Context, orgID uuid.UUID, rol string, clienteScope *uuid.UUID, cargaID uuid.UUID) ([]Evento, error) {
	if s.Carga != nil {
		if err := s.Carga.EnsureVisible(ctx, orgID, rol, clienteScope, cargaID); err != nil {
			return nil, httperr.NotFound("carga no encontrada")
		}
	}
	return s.Repo.ListByCarga(ctx, orgID, cargaID)
}
