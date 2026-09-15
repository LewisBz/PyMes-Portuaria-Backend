package organizacion

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

func (p *Postgres) GetByID(ctx context.Context, id uuid.UUID) (*Organizacion, error) {
	var o Organizacion
	err := db.WithContext(ctx, p.DB).First(&o, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, httperr.NotFound("organización no encontrada")
	}
	return &o, err
}

func (p *Postgres) IsActiva(ctx context.Context, id uuid.UUID) (bool, error) {
	o, err := p.GetByID(ctx, id)
	if err != nil {
		return false, err
	}
	return o.Activa, nil
}
