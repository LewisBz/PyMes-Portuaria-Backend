package usuarios

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, orgID, id uuid.UUID) (*Usuario, error)
	GetByEmail(ctx context.Context, orgID uuid.UUID, email string) (*Usuario, error)
	GetByEmailAnyOrg(ctx context.Context, email string) (*Usuario, error)
	List(ctx context.Context, orgID uuid.UUID) ([]Usuario, error)
	Create(ctx context.Context, u *Usuario) error
	Update(ctx context.Context, u *Usuario) error
}
