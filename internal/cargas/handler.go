package cargas

import (
	"net/http"
	"time"

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
	var est *Estado
	if q := c.Query("estado"); q != "" {
		e := Estado(q)
		est = &e
	}
	var cid *uuid.UUID
	if q := c.Query("cliente_id"); q != "" {
		id, err := uuid.Parse(q)
		if err == nil {
			cid = &id
		}
	}
	list, err := h.Svc.List(c.Request.Context(), cl.OrganizacionID, cl.Rol, cl.ClienteID, est, cid)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	out := make([]DTO, 0, len(list))
	for _, x := range list {
		out = append(out, x.ToDTO())
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) Create(c *gin.Context) {
	cl, err := httpx.MustClaims(c)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	var req struct {
		ClienteID         string  `json:"cliente_id"`
		Referencia        string  `json:"referencia"`
		Descripcion       *string `json:"descripcion"`
		Origen            *string `json:"origen"`
		Destino           *string `json:"destino"`
		FechaComprometida *string `json:"fecha_comprometida"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Abort(c, httperr.BadRequest("solicitud inválida"))
		return
	}
	cid, err := uuid.Parse(req.ClienteID)
	if err != nil {
		httpx.Abort(c, httperr.BadRequest("cliente_id inválido"))
		return
	}
	var fc *time.Time
	if req.FechaComprometida != nil && *req.FechaComprometida != "" {
		t, err := time.Parse("2006-01-02", *req.FechaComprometida)
		if err != nil {
			httpx.Abort(c, httperr.BadRequest("fecha_comprometida inválida"))
			return
		}
		fc = &t
	}
	cg, err := h.Svc.Create(c.Request.Context(), cl.OrganizacionID, cl.UserID, cid, req.Referencia, req.Descripcion, req.Origen, req.Destino, fc)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	c.JSON(http.StatusCreated, cg.ToDTO())
}

func (h *Handler) Get(c *gin.Context) {
	cl, err := httpx.MustClaims(c)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Abort(c, httperr.NotFound("carga no encontrada"))
		return
	}
	cg, err := h.Svc.Get(c.Request.Context(), cl.OrganizacionID, cl.Rol, cl.ClienteID, id)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	c.JSON(http.StatusOK, cg.ToDTO())
}

func (h *Handler) CambiarEstado(c *gin.Context) {
	cl, err := httpx.MustClaims(c)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Abort(c, httperr.NotFound("carga no encontrada"))
		return
	}
	var req struct {
		Estado Estado `json:"estado"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Estado == "" {
		httpx.Abort(c, httperr.BadRequest("estado requerido"))
		return
	}
	cg, err := h.Svc.CambiarEstado(c.Request.Context(), cl.OrganizacionID, cl.UserID, id, req.Estado)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	c.JSON(http.StatusOK, cg.ToDTO())
}
