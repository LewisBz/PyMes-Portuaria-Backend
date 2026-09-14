package usuarios

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
	list, err := h.Svc.List(c.Request.Context(), cl.OrganizacionID)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	out := make([]Public, 0, len(list))
	for _, u := range list {
		out = append(out, u.Public())
	}
	c.JSON(http.StatusOK, out)
}

type createReq struct {
	Email     string  `json:"email" binding:"required,email"`
	Password  string  `json:"password" binding:"required"`
	Nombre    string  `json:"nombre" binding:"required"`
	Rol       Rol     `json:"rol" binding:"required"`
	ClienteID *string `json:"cliente_id"`
}

func (h *Handler) Create(c *gin.Context) {
	cl, err := httpx.MustClaims(c)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	var req createReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Abort(c, httperr.BadRequest("solicitud inválida"))
		return
	}
	var cid *uuid.UUID
	if req.ClienteID != nil && *req.ClienteID != "" {
		id, err := uuid.Parse(*req.ClienteID)
		if err != nil {
			httpx.Abort(c, httperr.BadRequest("cliente_id inválido"))
			return
		}
		cid = &id
	}
	u, err := h.Svc.Create(c.Request.Context(), cl.OrganizacionID, req.Email, req.Password, req.Nombre, req.Rol, cid)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	c.JSON(http.StatusCreated, u.Public())
}

type patchReq struct {
	Nombre *string `json:"nombre"`
	Activo *bool   `json:"activo"`
	Rol    *Rol    `json:"rol"`
}

func (h *Handler) Patch(c *gin.Context) {
	cl, err := httpx.MustClaims(c)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Abort(c, httperr.BadRequest("solicitud inválida"))
		return
	}
	var req patchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Abort(c, httperr.BadRequest("solicitud inválida"))
		return
	}
	u, err := h.Svc.Patch(c.Request.Context(), cl.OrganizacionID, id, req.Nombre, req.Activo, req.Rol)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	c.JSON(http.StatusOK, u.Public())
}
