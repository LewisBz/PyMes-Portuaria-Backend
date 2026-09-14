package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pymes-portuaria/backend/internal/organizacion"
	"github.com/pymes-portuaria/backend/internal/platform/httpx"
	"github.com/pymes-portuaria/backend/internal/usuarios"
	"github.com/pymes-portuaria/backend/pkg/httperr"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Users  usuarios.Repository
	Orgs   organizacion.Repository
	Secret string
	Expiry time.Duration
}

type LoginResult struct {
	AccessToken string          `json:"access_token"`
	TokenType   string          `json:"token_type"`
	Usuario     usuarios.Public `json:"usuario"`
}

func (s *Service) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	u, err := s.Users.GetByEmailAnyOrg(ctx, email)
	if err != nil {
		return nil, httperr.Unauthorized("credenciales inválidas")
	}
	if !u.Activo {
		return nil, httperr.Unauthorized("credenciales inválidas")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, httperr.Unauthorized("credenciales inválidas")
	}
	org, err := s.Orgs.GetByID(ctx, u.OrganizacionID)
	if err != nil || !org.Activa {
		return nil, httperr.Unauthorized("credenciales inválidas")
	}
	claims := httpx.TokenClaims{
		UserID:         u.ID.String(),
		OrganizacionID: u.OrganizacionID.String(),
		Rol:            string(u.Rol),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.Expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString([]byte(s.Secret))
	if err != nil {
		return nil, httperr.Internal("no se pudo emitir token")
	}
	return &LoginResult{
		AccessToken: signed,
		TokenType:   "Bearer",
		Usuario:     u.Public(),
	}, nil
}

func (s *Service) Me(ctx context.Context, orgID, userID uuid.UUID) (*usuarios.Usuario, error) {
	return s.Users.GetByID(ctx, orgID, userID)
}
