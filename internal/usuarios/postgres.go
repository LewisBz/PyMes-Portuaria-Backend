package usuarios

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

func (p *Postgres) GetByID(ctx context.Context, orgID, id uuid.UUID) (*Usuario, error) {
	var u Usuario
	err := db.WithContext(ctx, p.DB).Where("id = ? AND organizacion_id = ?", id, orgID).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, httperr.NotFound("usuario no encontrado")
	}
	return &u, err
}

func (p *Postgres) GetByEmail(ctx context.Context, orgID uuid.UUID, email string) (*Usuario, error) {
	var u Usuario
	err := db.WithContext(ctx, p.DB).Where("organizacion_id = ? AND email = ?", orgID, email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, httperr.NotFound("usuario no encontrado")
	}
	return &u, err
}

func (p *Postgres) GetByEmailAnyOrg(ctx context.Context, email string) (*Usuario, error) {
	var u Usuario
	err := db.WithContext(ctx, p.DB).Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, httperr.Unauthorized("credenciales inválidas")
	}
	return &u, err
}

func (p *Postgres) List(ctx context.Context, orgID uuid.UUID) ([]Usuario, error) {
	var list []Usuario
	err := db.WithContext(ctx, p.DB).Where("organizacion_id = ?", orgID).Order("email").Find(&list).Error
	return list, err
}

func (p *Postgres) Create(ctx context.Context, u *Usuario) error {
	return db.WithContext(ctx, p.DB).Create(u).Error
}

func (p *Postgres) Update(ctx context.Context, u *Usuario) error {
	return db.WithContext(ctx, p.DB).Save(u).Error
}
