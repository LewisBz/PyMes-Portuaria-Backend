package cargas

import (
	"time"

	"github.com/google/uuid"
)

type Estado string

const (
	EstadoRegistrada Estado = "registrada"
	EstadoEnTransito Estado = "en_transito"
	EstadoEnPuerto   Estado = "en_puerto"
	EstadoEntregada  Estado = "entregada"
	EstadoCancelada  Estado = "cancelada"
)

func (e Estado) Terminal() bool {
	return e == EstadoEntregada || e == EstadoCancelada
}

type Carga struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrganizacionID    uuid.UUID `gorm:"type:uuid"`
	ClienteID         uuid.UUID `gorm:"type:uuid"`
	Referencia        string
	Descripcion       *string
	Origen            *string
	Destino           *string
	Estado            Estado
	FechaComprometida *time.Time `gorm:"type:date"`
	FechaEntrega      *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (Carga) TableName() string { return "cargas" }

type DTO struct {
	ID                uuid.UUID  `json:"id"`
	ClienteID         uuid.UUID  `json:"cliente_id"`
	Referencia        string     `json:"referencia"`
	Descripcion       *string    `json:"descripcion"`
	Origen            *string    `json:"origen"`
	Destino           *string    `json:"destino"`
	Estado            Estado     `json:"estado"`
	Retrasada         bool       `json:"retrasada"`
	FechaComprometida *string    `json:"fecha_comprometida"`
	FechaEntrega      *time.Time `json:"fecha_entrega"`
	CreatedAt         time.Time  `json:"created_at"`
}

func (c Carga) ToDTO() DTO {
	retrasada := false
	if c.FechaComprometida != nil && c.Estado != EstadoEntregada {
		retrasada = c.FechaComprometida.Before(time.Now())
	}
	var fc *string
	if c.FechaComprometida != nil {
		s := c.FechaComprometida.Format("2006-01-02")
		fc = &s
	}
	return DTO{
		ID: c.ID, ClienteID: c.ClienteID, Referencia: c.Referencia,
		Descripcion: c.Descripcion, Origen: c.Origen, Destino: c.Destino,
		Estado: c.Estado, Retrasada: retrasada, FechaComprometida: fc,
		FechaEntrega: c.FechaEntrega, CreatedAt: c.CreatedAt,
	}
}
