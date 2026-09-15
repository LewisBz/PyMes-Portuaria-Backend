package cargas

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pymes-portuaria/backend/internal/clientes"
	"github.com/pymes-portuaria/backend/internal/trazabilidad"
	"github.com/pymes-portuaria/backend/pkg/httperr"
	"gorm.io/gorm"
)

type Service struct {
	Repo     Repository
	Clientes clientes.Repository
	Events   trazabilidad.Appender
	DB       *gorm.DB
}

func (s *Service) visibleCliente(clCliente *uuid.UUID, cargaCliente uuid.UUID) error {
	if clCliente == nil {
		return nil
	}
	if *clCliente != cargaCliente {
		return httperr.NotFound("carga no encontrada")
	}
	return nil
}

func (s *Service) Create(ctx context.Context, orgID, actorID uuid.UUID, clienteID uuid.UUID, ref string, desc, origen, destino *string, fechaComp *time.Time) (*Carga, error) {
	if ref == "" {
		return nil, httperr.BadRequest("referencia requerida")
	}
	if _, err := s.Clientes.GetByID(ctx, orgID, clienteID); err != nil {
		return nil, httperr.BadRequest("cliente_id inválido")
	}
	c := &Carga{
		ID: uuid.New(), OrganizacionID: orgID, ClienteID: clienteID,
		Referencia: ref, Descripcion: desc, Origen: origen, Destino: destino,
		Estado: EstadoRegistrada, FechaComprometida: fechaComp,
	}
	if err := s.Repo.Create(ctx, c); err != nil {
		return nil, httperr.Conflict("no se pudo crear la carga")
	}
	if s.Events != nil {
		_ = s.Events.Append(ctx, trazabilidad.Evento{
			ID: uuid.New(), OrganizacionID: orgID, CargaID: c.ID,
			Tipo: trazabilidad.TipoCargaCreada, ActorID: actorID,
			OccurredAt: time.Now(), EstadoNuevo: strPtr(string(c.Estado)),
		})
	}
	return c, nil
}

func (s *Service) List(ctx context.Context, orgID uuid.UUID, rol string, clienteScope *uuid.UUID, estado *Estado, filterCliente *uuid.UUID) ([]Carga, error) {
	var scope *uuid.UUID
	if rol == "cliente" {
		scope = clienteScope
	} else if filterCliente != nil {
		scope = filterCliente
	}
	return s.Repo.List(ctx, orgID, scope, estado)
}

func (s *Service) Get(ctx context.Context, orgID uuid.UUID, rol string, clienteScope *uuid.UUID, id uuid.UUID) (*Carga, error) {
	c, err := s.Repo.GetByID(ctx, orgID, id)
	if err != nil {
		return nil, httperr.NotFound("carga no encontrada")
	}
	if err := s.visibleCliente(clienteScopeIf(rol, clienteScope), c.ClienteID); err != nil {
		return nil, err
	}
	return c, nil
}

func clienteScopeIf(rol string, id *uuid.UUID) *uuid.UUID {
	if rol == "cliente" {
		return id
	}
	return nil
}

func (s *Service) CambiarEstado(ctx context.Context, orgID, actorID uuid.UUID, cargaID uuid.UUID, nuevo Estado) (*Carga, error) {
	if s.DB == nil {
		return s.cambiar(ctx, orgID, actorID, cargaID, nuevo)
	}
	var out *Carga
	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		origRepo, origEv := s.Repo, s.Events
		s.Repo = &Postgres{DB: tx}
		if s.Events != nil {
			s.Events = &trazabilidad.Postgres{DB: tx}
		}
		defer func() { s.Repo, s.Events = origRepo, origEv }()
		c, err := s.cambiar(ctx, orgID, actorID, cargaID, nuevo)
		out = c
		return err
	})
	return out, err
}

func (s *Service) cambiar(ctx context.Context, orgID, actorID, cargaID uuid.UUID, nuevo Estado) (*Carga, error) {
	c, err := s.Repo.GetByID(ctx, orgID, cargaID)
	if err != nil {
		return nil, httperr.NotFound("carga no encontrada")
	}
	if c.Estado == nuevo {
		return nil, httperr.Conflict("el estado no cambió")
	}
	if c.Estado.Terminal() {
		return nil, httperr.Conflict("carga en estado terminal")
	}
	prev := c.Estado
	c.Estado = nuevo
	now := time.Now()
	if nuevo == EstadoEntregada {
		c.FechaEntrega = &now
	}
	if err := s.Repo.Update(ctx, c); err != nil {
		return nil, err
	}
	if s.Events != nil {
		ps := string(prev)
		ns := string(nuevo)
		if err := s.Events.Append(ctx, trazabilidad.Evento{
			ID: uuid.New(), OrganizacionID: orgID, CargaID: c.ID,
			Tipo: trazabilidad.TipoEstadoCambiado, ActorID: actorID,
			OccurredAt: now, EstadoAnterior: &ps, EstadoNuevo: &ns,
		}); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (s *Service) EnsureVisible(ctx context.Context, orgID uuid.UUID, rol string, clienteScope *uuid.UUID, cargaID uuid.UUID) error {
	_, err := s.Get(ctx, orgID, rol, clienteScope, cargaID)
	return err
}

func strPtr(s string) *string { return &s }
