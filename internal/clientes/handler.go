package clientes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pymes-portuaria/backend/internal/platform/httpx"
	"github.com/pymes-portuaria/backend/pkg/httperr"
)

type Handler struct {
	Svc *Service
}

func (h *Handler) List(c *gin.Context) {
	cl, err := httpx.MustClaims(c)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	list, err := h.Svc.List(c.Request.Context(), cl.OrganizacionID)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) Create(c *gin.Context) {
	cl, err := httpx.MustClaims(c)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	var req struct {
		RazonSocial string  `json:"razon_social"`
		NIT         *string `json:"nit"`
		Email       *string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Abort(c, httperr.BadRequest("solicitud inválida"))
		return
	}
	cli, err := h.Svc.Create(c.Request.Context(), cl.OrganizacionID, req.RazonSocial, req.NIT, req.Email)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	c.JSON(http.StatusCreated, cli)
}
