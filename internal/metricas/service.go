package metricas

import (
	"context"

	"github.com/google/uuid"
	"github.com/pymes-portuaria/backend/internal/platform/db"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func (s *Service) Get(ctx context.Context, orgID uuid.UUID) (Snapshot, error) {
	gdb := db.WithContext(ctx, s.DB)
	var snap Snapshot
	if err := gdb.Raw(`SELECT COUNT(*) FROM cargas WHERE organizacion_id = ? AND estado NOT IN ('entregada','cancelada')`, orgID).Scan(&snap.CargasActivas).Error; err != nil {
		return snap, err
	}
	if err := gdb.Raw(`SELECT COUNT(*) FROM cargas WHERE organizacion_id = ? AND fecha_comprometida IS NOT NULL AND fecha_comprometida < NOW() AND estado <> 'entregada'`, orgID).Scan(&snap.CargasRetrasadas).Error; err != nil {
		return snap, err
	}
	if err := gdb.Raw(`SELECT COUNT(*) FROM incidencias WHERE organizacion_id = ? AND estado = 'abierta'`, orgID).Scan(&snap.IncidenciasAbiertas).Error; err != nil {
		return snap, err
	}
	var avg *float64
	if err := gdb.Raw(`SELECT AVG(EXTRACT(EPOCH FROM (fecha_entrega - created_at))/3600.0) FROM cargas WHERE organizacion_id = ? AND estado = 'entregada' AND fecha_entrega IS NOT NULL`, orgID).Row().Scan(&avg); err != nil {
		snap.TiempoPromedioTotalEntregaHoras = nil
		return snap, nil
	}
	snap.TiempoPromedioTotalEntregaHoras = avg
	return snap, nil
}
