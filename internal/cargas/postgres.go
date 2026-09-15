package cargas

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

func (p *Postgres) Create(ctx context.Context, c *Carga) error {
	return db.WithContext(ctx, p.DB).Create(c).Error
}

func (p *Postgres) List(ctx context.Context, orgID uuid.UUID, clienteID *uuid.UUID, estado *Estado) ([]Carga, error) {
	q := db.WithContext(ctx, p.DB).Where("organizacion_id = ?", orgID)
	if clienteID != nil {
		q = q.Where("cliente_id = ?", *clienteID)
	}
	if estado != nil {
		q = q.Where("estado = ?", *estado)
	}
	var list []Carga
	err := q.Order("created_at desc").Find(&list).Error
	return list, err
}

func (p *Postgres) GetByID(ctx context.Context, orgID, id uuid.UUID) (*Carga, error) {
	var c Carga
	err := db.WithContext(ctx, p.DB).Where("id = ? AND organizacion_id = ?", id, orgID).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, httperr.NotFound("carga no encontrada")
	}
	return &c, err
}

func (p *Postgres) Update(ctx context.Context, c *Carga) error {
	return db.WithContext(ctx, p.DB).Save(c).Error
}
