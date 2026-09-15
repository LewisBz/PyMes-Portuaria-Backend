package app

import (
	"github.com/gin-gonic/gin"
	"github.com/pymes-portuaria/backend/internal/platform/httpx"
	"github.com/pymes-portuaria/backend/internal/usuarios"
	"github.com/pymes-portuaria/backend/pkg/httperr"
)

func RejectInactiveUser(users usuarios.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		cl, err := httpx.MustClaims(c)
		if err != nil {
			httpx.Abort(c, err)
			return
		}
		u, err := users.GetByID(c.Request.Context(), cl.OrganizacionID, cl.UserID)
		if err != nil || !u.Activo {
			httpx.Abort(c, httperr.Unauthorized("usuario inactivo"))
			return
		}
		c.Next()
	}
}

func EnrichCliente(users usuarios.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		cl, ok := httpx.GetClaims(c)
		if !ok {
			c.Next()
			return
		}
		u, err := users.GetByID(c.Request.Context(), cl.OrganizacionID, cl.UserID)
		if err == nil {
			cl.ClienteID = u.ClienteID
			httpx.SetClaims(c, cl)
		}
		c.Next()
	}
}
