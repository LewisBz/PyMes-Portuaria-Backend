package trazabilidad

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Tipo string

const (
	TipoCargaCreada       Tipo = "carga_creada"
	TipoEstadoCambiado    Tipo = "estado_cambiado"
	TipoIncidenciaAbierta Tipo = "incidencia_abierta"
	TipoIncidenciaCerrada Tipo = "incidencia_cerrada"
	TipoDocumentoAdjunto  Tipo = "documento_adjunto"
)

type Evento struct {
	ID             uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	OrganizacionID uuid.UUID       `gorm:"type:uuid" json:"-"`
	CargaID        uuid.UUID       `gorm:"type:uuid" json:"carga_id"`
	Tipo           Tipo            `json:"tipo"`
	ActorID        uuid.UUID       `gorm:"type:uuid" json:"actor_id"`
	OccurredAt     time.Time       `json:"occurred_at"`
	EstadoAnterior *string         `json:"estado_anterior"`
	EstadoNuevo    *string         `json:"estado_nuevo"`
	Payload        json.RawMessage `gorm:"type:jsonb" json:"payload"`
}

func (Evento) TableName() string { return "eventos_trazabilidad" }
