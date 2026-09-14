package documentos

import (
	"io"
	"net/http"
	"strconv"

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

func (h *Handler) Upload(c *gin.Context) {
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
	fh, err := c.FormFile("file")
	if err != nil {
		httpx.Abort(c, httperr.BadRequest("file requerido"))
		return
	}
	f, err := fh.Open()
	if err != nil {
		httpx.Abort(c, httperr.BadRequest("file inválido"))
		return
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		httpx.Abort(c, httperr.BadRequest("file inválido"))
		return
	}
	ctype := fh.Header.Get("Content-Type")
	if ctype == "" || ctype == "application/octet-stream" {
		ctype = http.DetectContentType(data)
	}
	visible, _ := strconv.ParseBool(c.PostForm("visible_cliente"))
	d, err := h.Svc.Upload(c.Request.Context(), cl.OrganizacionID, cl.UserID, cid, fh.Filename, ctype, data, visible)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	c.JSON(http.StatusCreated, d)
}

func (h *Handler) Download(c *gin.Context) {
	cl, err := httpx.MustClaims(c)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Abort(c, httperr.NotFound("documento no encontrado"))
		return
	}
	d, data, err := h.Svc.Download(c.Request.Context(), cl.OrganizacionID, cl.Rol, cl.ClienteID, id)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	c.Data(http.StatusOK, d.ContentType, data)
}
