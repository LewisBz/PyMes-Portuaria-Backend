package documentos

import (
	"time"

	"github.com/google/uuid"
)

type Documento struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	OrganizacionID uuid.UUID `gorm:"type:uuid" json:"-"`
	CargaID        uuid.UUID `gorm:"type:uuid" json:"carga_id"`
	NombreOriginal string    `json:"nombre_original"`
	ContentType    string    `json:"content_type"`
	TamanoBytes    int       `json:"tamano_bytes"`
	StorageKey     string    `json:"-"`
	VisibleCliente bool      `json:"visible_cliente"`
	SubidoPor      uuid.UUID `gorm:"type:uuid" json:"-"`
	CreatedAt      time.Time `json:"created_at"`
}

func (Documento) TableName() string { return "documentos" }
