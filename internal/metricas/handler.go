package metricas

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pymes-portuaria/backend/internal/platform/httpx"
)

type Getter interface {
	Get(ctx context.Context, orgID uuid.UUID) (Snapshot, error)
}

type Handler struct {
	Svc Getter
}

func (h *Handler) Get(c *gin.Context) {
	cl, err := httpx.MustClaims(c)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	s, err := h.Svc.Get(c.Request.Context(), cl.OrganizacionID)
	if err != nil {
		httpx.Abort(c, err)
		return
	}
	c.JSON(http.StatusOK, s)
}
