package clientes

import (
	"context"

	"github.com/google/uuid"
	"github.com/pymes-portuaria/backend/pkg/httperr"
)

type Service struct {
	Repo Repository
}

func (s *Service) Create(ctx context.Context, orgID uuid.UUID, razon string, nit, email *string) (*Cliente, error) {
	if razon == "" {
		return nil, httperr.BadRequest("razon_social requerida")
	}
	c := &Cliente{ID: uuid.New(), OrganizacionID: orgID, RazonSocial: razon, NIT: nit, Email: email, Activo: true}
	if err := s.Repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) List(ctx context.Context, orgID uuid.UUID) ([]Cliente, error) {
	return s.Repo.List(ctx, orgID)
}
