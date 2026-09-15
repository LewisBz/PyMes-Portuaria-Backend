package documentos

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/pymes-portuaria/backend/internal/platform/db"
	"github.com/pymes-portuaria/backend/pkg/httperr"
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func (p *Postgres) Create(ctx context.Context, d *Documento) error {
	return db.WithContext(ctx, p.DB).Create(d).Error
}

func (p *Postgres) ListByCarga(ctx context.Context, orgID, cargaID uuid.UUID, onlyVisible bool) ([]Documento, error) {
	q := db.WithContext(ctx, p.DB).Where("organizacion_id = ? AND carga_id = ?", orgID, cargaID)
	if onlyVisible {
		q = q.Where("visible_cliente = true")
	}
	var list []Documento
	err := q.Find(&list).Error
	return list, err
}

func (p *Postgres) GetByID(ctx context.Context, orgID, id uuid.UUID) (*Documento, error) {
	var d Documento
	err := db.WithContext(ctx, p.DB).Where("id = ? AND organizacion_id = ?", id, orgID).First(&d).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, httperr.NotFound("documento no encontrado")
	}
	return &d, err
}
