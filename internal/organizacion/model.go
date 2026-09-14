package organizacion

import (
	"time"

	"github.com/google/uuid"
)

type Organizacion struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Nombre    string
	Activa    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Organizacion) TableName() string { return "organizaciones" }
