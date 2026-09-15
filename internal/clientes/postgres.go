package clientes

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

func (p *Postgres) Create(ctx context.Context, c *Cliente) error {
	return db.WithContext(ctx, p.DB).Create(c).Error
}

func (p *Postgres) List(ctx context.Context, orgID uuid.UUID) ([]Cliente, error) {
	var list []Cliente
	err := db.WithContext(ctx, p.DB).Where("organizacion_id = ?", orgID).Find(&list).Error
	return list, err
}

func (p *Postgres) GetByID(ctx context.Context, orgID, id uuid.UUID) (*Cliente, error) {
	var c Cliente
	err := db.WithContext(ctx, p.DB).Where("id = ? AND organizacion_id = ?", id, orgID).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, httperr.NotFound("cliente no encontrado")
	}
	return &c, err
}
