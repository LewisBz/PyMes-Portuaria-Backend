package usuarios

import (
	"time"

	"github.com/google/uuid"
)

type Rol string

const (
	RolAdministrador Rol = "administrador"
	RolOperador      Rol = "operador"
	RolGerente       Rol = "gerente"
	RolCliente       Rol = "cliente"
)

type Usuario struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrganizacionID uuid.UUID `gorm:"type:uuid"`
	Email          string
	PasswordHash   string
	Nombre         string
	Rol            Rol
	ClienteID      *uuid.UUID `gorm:"type:uuid"`
	Activo         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (Usuario) TableName() string { return "usuarios" }

type Public struct {
	ID             uuid.UUID  `json:"id"`
	OrganizacionID uuid.UUID  `json:"organizacion_id"`
	Email          string     `json:"email"`
	Nombre         string     `json:"nombre"`
	Rol            Rol        `json:"rol"`
	ClienteID      *uuid.UUID `json:"cliente_id"`
	Activo         bool       `json:"activo"`
}

func (u Usuario) Public() Public {
	return Public{
		ID: u.ID, OrganizacionID: u.OrganizacionID, Email: u.Email,
		Nombre: u.Nombre, Rol: u.Rol, ClienteID: u.ClienteID, Activo: u.Activo,
	}
}
