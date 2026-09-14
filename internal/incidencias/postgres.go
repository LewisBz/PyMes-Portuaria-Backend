package incidencias

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

func (p *Postgres) Create(ctx context.Context, i *Incidencia) error {
	return db.WithContext(ctx, p.DB).Create(i).Error
}

func (p *Postgres) ListByCarga(ctx context.Context, orgID, cargaID uuid.UUID) ([]Incidencia, error) {
	var list []Incidencia
	err := db.WithContext(ctx, p.DB).Where("organizacion_id = ? AND carga_id = ?", orgID, cargaID).Find(&list).Error
	return list, err
}

func (p *Postgres) GetByID(ctx context.Context, orgID, id uuid.UUID) (*Incidencia, error) {
	var i Incidencia
	err := db.WithContext(ctx, p.DB).Where("id = ? AND organizacion_id = ?", id, orgID).First(&i).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, httperr.NotFound("incidencia no encontrada")
	}
	return &i, err
}

func (p *Postgres) Update(ctx context.Context, i *Incidencia) error {
	return db.WithContext(ctx, p.DB).Save(i).Error
}
