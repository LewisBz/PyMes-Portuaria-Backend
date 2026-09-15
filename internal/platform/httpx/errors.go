package httpx

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pymes-portuaria/backend/pkg/httperr"
)

func Abort(c *gin.Context, err error) {
	var he *httperr.Error
	if errors.As(err, &he) {
		c.AbortWithStatusJSON(he.Code, he)
		return
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, httperr.Internal("error interno"))
}

func JSONError(c *gin.Context, err error) {
	Abort(c, err)
}
