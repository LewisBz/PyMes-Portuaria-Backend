package incidencias

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pymes-portuaria/backend/internal/cargas"
	"github.com/pymes-portuaria/backend/internal/trazabilidad"
	"github.com/pymes-portuaria/backend/pkg/httperr"
	"gorm.io/gorm"
)

type Service struct {
	Repo   Repository
	Cargas *cargas.Service
	Events trazabilidad.Appender
	DB     *gorm.DB
}

func (s *Service) Open(ctx context.Context, orgID, actorID, cargaID uuid.UUID, titulo, desc, tipo string) (*Incidencia, error) {
	if s.DB == nil {
		return s.open(ctx, orgID, actorID, cargaID, titulo, desc, tipo)
	}
	var out *Incidencia
	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		origRepo, origEv := s.Repo, s.Events
		s.Repo = &Postgres{DB: tx}
		if s.Events != nil {
			s.Events = &trazabilidad.Postgres{DB: tx}
		}
		defer func() { s.Repo, s.Events = origRepo, origEv }()
		i, err := s.open(ctx, orgID, actorID, cargaID, titulo, desc, tipo)
		out = i
		return err
	})
	return out, err
}

func (s *Service) open(ctx context.Context, orgID, actorID, cargaID uuid.UUID, titulo, desc, tipo string) (*Incidencia, error) {
	if _, err := s.Cargas.Get(ctx, orgID, "operador", nil, cargaID); err != nil {
		return nil, err
	}
	if tipo == "" {
		tipo = "otro"
	}
	i := &Incidencia{
		ID: uuid.New(), OrganizacionID: orgID, CargaID: cargaID,
		Titulo: titulo, Descripcion: desc, Tipo: tipo, Estado: EstadoAbierta, AbiertaPor: actorID,
	}
	if err := s.Repo.Create(ctx, i); err != nil {
		return nil, err
	}
	if s.Events != nil {
		if err := s.Events.Append(ctx, trazabilidad.Evento{
			ID: uuid.New(), OrganizacionID: orgID, CargaID: cargaID,
			Tipo: trazabilidad.TipoIncidenciaAbierta, ActorID: actorID, OccurredAt: time.Now(),
		}); err != nil {
			return nil, err
		}
	}
	return i, nil
}

func (s *Service) List(ctx context.Context, orgID uuid.UUID, rol string, clienteScope *uuid.UUID, cargaID uuid.UUID) ([]Incidencia, error) {
	if _, err := s.Cargas.Get(ctx, orgID, rol, clienteScope, cargaID); err != nil {
		return nil, err
	}
	return s.Repo.ListByCarga(ctx, orgID, cargaID)
}

func (s *Service) Cerrar(ctx context.Context, orgID, actorID, id uuid.UUID) (*Incidencia, error) {
	if s.DB == nil {
		return s.cerrar(ctx, orgID, actorID, id)
	}
	var out *Incidencia
	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		origRepo, origEv := s.Repo, s.Events
		s.Repo = &Postgres{DB: tx}
		if s.Events != nil {
			s.Events = &trazabilidad.Postgres{DB: tx}
		}
		defer func() { s.Repo, s.Events = origRepo, origEv }()
		i, err := s.cerrar(ctx, orgID, actorID, id)
		out = i
		return err
	})
	return out, err
}

func (s *Service) cerrar(ctx context.Context, orgID, actorID, id uuid.UUID) (*Incidencia, error) {
	i, err := s.Repo.GetByID(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if i.Estado == EstadoCerrada {
		return nil, httperr.Conflict("incidencia ya cerrada")
	}
	now := time.Now()
	i.Estado = EstadoCerrada
	i.CerradaPor = &actorID
	i.ClosedAt = &now
	if err := s.Repo.Update(ctx, i); err != nil {
		return nil, err
	}
	if s.Events != nil {
		if err := s.Events.Append(ctx, trazabilidad.Evento{
			ID: uuid.New(), OrganizacionID: orgID, CargaID: i.CargaID,
			Tipo: trazabilidad.TipoIncidenciaCerrada, ActorID: actorID, OccurredAt: now,
		}); err != nil {
			return nil, err
		}
	}
	return i, nil
}
