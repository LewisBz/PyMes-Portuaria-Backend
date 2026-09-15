package incidencias

import (
	"time"

	"github.com/google/uuid"
)

type Estado string

const (
	EstadoAbierta Estado = "abierta"
	EstadoCerrada Estado = "cerrada"
)

type Incidencia struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	OrganizacionID uuid.UUID  `gorm:"type:uuid" json:"-"`
	CargaID        uuid.UUID  `gorm:"type:uuid" json:"carga_id"`
	Titulo         string     `json:"titulo"`
	Descripcion    string     `json:"descripcion"`
	Tipo           string     `json:"tipo"`
	Estado         Estado     `json:"estado"`
	AbiertaPor     uuid.UUID  `gorm:"type:uuid" json:"-"`
	CerradaPor     *uuid.UUID `gorm:"type:uuid" json:"-"`
	CreatedAt      time.Time  `json:"created_at"`
	ClosedAt       *time.Time `json:"closed_at"`
}

func (Incidencia) TableName() string { return "incidencias" }
