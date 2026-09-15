package trazabilidad

import (
	"context"

	"github.com/google/uuid"
	"github.com/pymes-portuaria/backend/internal/platform/db"
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func (p *Postgres) Append(ctx context.Context, e Evento) error {
	if e.Payload == nil {
		e.Payload = []byte("{}")
	}
	return db.WithContext(ctx, p.DB).Create(&e).Error
}

func (p *Postgres) ListByCarga(ctx context.Context, orgID, cargaID uuid.UUID) ([]Evento, error) {
	var list []Evento
	err := db.WithContext(ctx, p.DB).
		Where("organizacion_id = ? AND carga_id = ?", orgID, cargaID).
		Order("occurred_at asc").
		Find(&list).Error
	return list, err
}
