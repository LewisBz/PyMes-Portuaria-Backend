package documentos

import (
	"bytes"
	"context"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/pymes-portuaria/backend/internal/cargas"
	"github.com/pymes-portuaria/backend/internal/platform/storage"
	"github.com/pymes-portuaria/backend/internal/trazabilidad"
	"github.com/pymes-portuaria/backend/pkg/httperr"
)

const maxBytes = 10 * 1024 * 1024

var allowed = map[string]struct{}{
	"application/pdf": {},
	"image/jpeg":      {},
	"image/png":       {},
}

type Service struct {
	Repo   Repository
	Store  storage.ObjectStore
	Cargas *cargas.Service
	Events trazabilidad.Appender
}

func (s *Service) Upload(ctx context.Context, orgID, actorID, cargaID uuid.UUID, name, ctype string, data []byte, visible bool) (*Documento, error) {
	if _, err := s.Cargas.Get(ctx, orgID, "operador", nil, cargaID); err != nil {
		return nil, err
	}
	ctype = sniffContentType(ctype, data)
	if _, ok := allowed[ctype]; !ok {
		return nil, httperr.BadRequest("tipo de archivo no permitido")
	}
	if len(data) == 0 || len(data) > maxBytes {
		return nil, httperr.BadRequest("tamaño de archivo inválido")
	}
	id := uuid.New()
	key := path.Join(orgID.String(), cargaID.String(), id.String())
	if err := s.Store.Put(ctx, key, data); err != nil {
		return nil, httperr.Internal("no se pudo guardar el archivo")
	}
	d := &Documento{
		ID: id, OrganizacionID: orgID, CargaID: cargaID,
		NombreOriginal: name, ContentType: ctype, TamanoBytes: len(data),
		StorageKey: key, VisibleCliente: visible, SubidoPor: actorID,
	}
	if err := s.Repo.Create(ctx, d); err != nil {
		return nil, err
	}
	if s.Events != nil {
		_ = s.Events.Append(ctx, trazabilidad.Evento{
			ID: uuid.New(), OrganizacionID: orgID, CargaID: cargaID,
			Tipo: trazabilidad.TipoDocumentoAdjunto, ActorID: actorID, OccurredAt: time.Now(),
		})
	}
	return d, nil
}

func (s *Service) List(ctx context.Context, orgID uuid.UUID, rol string, clienteScope *uuid.UUID, cargaID uuid.UUID) ([]Documento, error) {
	if _, err := s.Cargas.Get(ctx, orgID, rol, clienteScope, cargaID); err != nil {
		return nil, err
	}
	return s.Repo.ListByCarga(ctx, orgID, cargaID, rol == "cliente")
}

func (s *Service) Download(ctx context.Context, orgID uuid.UUID, rol string, clienteScope *uuid.UUID, id uuid.UUID) (*Documento, []byte, error) {
	d, err := s.Repo.GetByID(ctx, orgID, id)
	if err != nil {
		return nil, nil, err
	}
	if _, err := s.Cargas.Get(ctx, orgID, rol, clienteScope, d.CargaID); err != nil {
		return nil, nil, httperr.NotFound("documento no encontrado")
	}
	if rol == "cliente" && !d.VisibleCliente {
		return nil, nil, httperr.Forbidden("documento no autorizado")
	}
	b, err := s.Store.Get(ctx, d.StorageKey)
	if err != nil {
		return nil, nil, httperr.NotFound("documento no encontrado")
	}
	return d, b, nil
}

func sniffContentType(declared string, data []byte) string {
	switch {
	case bytes.HasPrefix(data, []byte("%PDF")):
		return "application/pdf"
	case bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47}):
		return "image/png"
	case bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF}):
		return "image/jpeg"
	default:
		return declared
	}
}
