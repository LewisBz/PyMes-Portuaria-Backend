package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pymes-portuaria/backend/internal/platform/httpx"
	"github.com/pymes-portuaria/backend/pkg/httperr"
)

type Handler struct {
	Svc *Service
}

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" || req.Password == "" {
		httpx.Abort(c, httperr.BadRequest("email y password requeridos"))
		return
	}
	res, err := h.Svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) Me(c *gin.Context) {
	cl, err := httpx.MustClaims(c)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	u, err := h.Svc.Me(c.Request.Context(), cl.OrganizacionID, cl.UserID)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	c.JSON(http.StatusOK, u.Public())
}
