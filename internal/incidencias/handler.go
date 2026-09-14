package incidencias

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	cid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Abort(c, httperr.NotFound("carga no encontrada"))
		return
	}
	list, err := h.Svc.List(c.Request.Context(), cl.OrganizacionID, cl.Rol, cl.ClienteID, cid)
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
	cid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Abort(c, httperr.NotFound("carga no encontrada"))
		return
	}
	var req struct {
		Titulo      string `json:"titulo"`
		Descripcion string `json:"descripcion"`
		Tipo        string `json:"tipo"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Titulo == "" || req.Descripcion == "" {
		httpx.Abort(c, httperr.BadRequest("titulo y descripcion requeridos"))
		return
	}
	i, err := h.Svc.Open(c.Request.Context(), cl.OrganizacionID, cl.UserID, cid, req.Titulo, req.Descripcion, req.Tipo)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	c.JSON(http.StatusCreated, i)
}

func (h *Handler) Cerrar(c *gin.Context) {
	cl, err := httpx.MustClaims(c)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Abort(c, httperr.NotFound("incidencia no encontrada"))
		return
	}
	i, err := h.Svc.Cerrar(c.Request.Context(), cl.OrganizacionID, cl.UserID, id)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	c.JSON(http.StatusOK, i)
}
