package httpx

import (
	"github.com/gin-gonic/gin"
	"github.com/pymes-portuaria/backend/pkg/httperr"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := map[string]struct{}{}
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *gin.Context) {
		cl, err := MustClaims(c)
		if err != nil {
			Abort(c, err)
			return
		}
		if _, ok := allowed[cl.Rol]; !ok {
			Abort(c, httperr.Forbidden("rol no autorizado"))
			return
		}
		c.Next()
	}
}
