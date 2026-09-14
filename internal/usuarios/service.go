package usuarios

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/pymes-portuaria/backend/pkg/httperr"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Repo Repository
	Cost int
}

func (s *Service) Create(ctx context.Context, orgID uuid.UUID, email, password, nombre string, rol Rol, clienteID *uuid.UUID) (*Usuario, error) {
	if s.Cost < 12 {
		s.Cost = 12
	}
	if len(password) < 8 {
		return nil, httperr.BadRequest("password mínimo 8 caracteres")
	}
	if rol == RolCliente && (clienteID == nil || *clienteID == uuid.Nil) {
		return nil, httperr.BadRequest("cliente_id requerido para rol cliente")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.Cost)
	if err != nil {
		return nil, httperr.Internal("no se pudo guardar la contraseña")
	}
	u := &Usuario{
		ID: uuid.New(), OrganizacionID: orgID, Email: strings.ToLower(email),
		PasswordHash: string(hash), Nombre: nombre, Rol: rol, ClienteID: clienteID, Activo: true,
	}
	if err := s.Repo.Create(ctx, u); err != nil {
		return nil, httperr.Conflict("no se pudo crear el usuario")
	}
	return u, nil
}

func (s *Service) List(ctx context.Context, orgID uuid.UUID) ([]Usuario, error) {
	return s.Repo.List(ctx, orgID)
}

func (s *Service) Patch(ctx context.Context, orgID, id uuid.UUID, nombre *string, activo *bool, rol *Rol) (*Usuario, error) {
	u, err := s.Repo.GetByID(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if nombre != nil {
		u.Nombre = *nombre
	}
	if activo != nil {
		u.Activo = *activo
	}
	if rol != nil {
		u.Rol = *rol
	}
	if err := s.Repo.Update(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}
