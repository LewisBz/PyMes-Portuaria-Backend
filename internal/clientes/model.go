package clientes

import (
	"time"

	"github.com/google/uuid"
)

type Cliente struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	OrganizacionID uuid.UUID `gorm:"type:uuid" json:"-"`
	RazonSocial    string    `json:"razon_social"`
	NIT            *string   `json:"nit"`
	Email          *string   `json:"email"`
	Activo         bool      `json:"activo"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

func (Cliente) TableName() string { return "clientes" }
