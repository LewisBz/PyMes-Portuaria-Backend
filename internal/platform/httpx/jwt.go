package httpx

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pymes-portuaria/backend/pkg/httperr"
)

type TokenClaims struct {
	UserID         string `json:"user_id"`
	OrganizacionID string `json:"organizacion_id"`
	Rol            string `json:"rol"`
	jwt.RegisteredClaims
}

type OrgActiveLookup func(ctx context.Context, orgID uuid.UUID) (bool, error)

func JWT(secret string, lookup OrgActiveLookup) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			Abort(c, httperr.Unauthorized("token requerido"))
			return
		}
		raw := strings.TrimPrefix(h, "Bearer ")
		tc := &TokenClaims{}
		tok, err := jwt.ParseWithClaims(raw, tc, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil || !tok.Valid {
			Abort(c, httperr.Unauthorized("token inválido"))
			return
		}
		uid, err1 := uuid.Parse(tc.UserID)
		oid, err2 := uuid.Parse(tc.OrganizacionID)
		if err1 != nil || err2 != nil || tc.Rol == "" {
			Abort(c, httperr.Unauthorized("token inválido"))
			return
		}
		if lookup != nil {
			ok, err := lookup(c.Request.Context(), oid)
			if err != nil {
				Abort(c, httperr.Internal("error de autorización"))
				return
			}
			if !ok {
				Abort(c, httperr.Unauthorized("organización inactiva"))
				return
			}
		}
		SetClaims(c, Claims{UserID: uid, OrganizacionID: oid, Rol: tc.Rol})
		c.Next()
	}
}
