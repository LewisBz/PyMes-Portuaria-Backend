package httpx

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pymes-portuaria/backend/pkg/httperr"
)

type Claims struct {
	UserID         uuid.UUID
	OrganizacionID uuid.UUID
	Rol            string
	ClienteID      *uuid.UUID
}

const claimsKey = "auth_claims"

func SetClaims(c *gin.Context, cl Claims) { c.Set(claimsKey, cl) }

func GetClaims(c *gin.Context) (Claims, bool) {
	v, ok := c.Get(claimsKey)
	if !ok {
		return Claims{}, false
	}
	cl, ok := v.(Claims)
	return cl, ok
}

func MustClaims(c *gin.Context) (Claims, error) {
	cl, ok := GetClaims(c)
	if !ok {
		return Claims{}, httperr.Unauthorized("no autenticado")
	}
	return cl, nil
}
