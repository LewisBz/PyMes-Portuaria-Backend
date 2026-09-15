package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pymes-portuaria/backend/internal/auth"
	"github.com/pymes-portuaria/backend/internal/cargas"
	"github.com/pymes-portuaria/backend/internal/clientes"
	"github.com/pymes-portuaria/backend/internal/documentos"
	"github.com/pymes-portuaria/backend/internal/incidencias"
	"github.com/pymes-portuaria/backend/internal/metricas"
	"github.com/pymes-portuaria/backend/internal/organizacion"
	"github.com/pymes-portuaria/backend/internal/platform/httpx"
	"github.com/pymes-portuaria/backend/internal/trazabilidad"
	"github.com/pymes-portuaria/backend/internal/usuarios"
)

type Deps struct {
	JWTSecret string
	Orgs      organizacion.Repository
	Users     usuarios.Repository
	Auth      *auth.Handler
	Usuarios  *usuarios.Handler
	Clientes  *clientes.Handler
	Cargas    *cargas.Handler
	Traz      *trazabilidad.Handler
	Incid     *incidencias.Handler
	Docs      *documentos.Handler
	Metricas  *metricas.Handler
}

func NewEngine(d Deps) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.MaxMultipartMemory = 11 << 20
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	v1 := r.Group("/api/v1")
	v1.POST("/auth/login", d.Auth.Login)

	authed := v1.Group("")
	authed.Use(
		httpx.JWT(d.JWTSecret, d.Orgs.IsActiva),
		RejectInactiveUser(d.Users),
		EnrichCliente(d.Users),
	)
	authed.GET("/auth/me", d.Auth.Me)

	admin := authed.Group("")
	admin.Use(httpx.RequireRole("administrador"))
	admin.GET("/usuarios", d.Usuarios.List)
	admin.POST("/usuarios", d.Usuarios.Create)
	admin.PATCH("/usuarios/:id", d.Usuarios.Patch)

	authed.GET("/clientes", httpx.RequireRole("administrador", "operador", "gerente"), d.Clientes.List)
	authed.POST("/clientes", httpx.RequireRole("operador", "administrador"), d.Clientes.Create)

	authed.GET("/cargas", d.Cargas.List)
	authed.POST("/cargas", httpx.RequireRole("operador", "administrador"), d.Cargas.Create)
	authed.GET("/cargas/:id", d.Cargas.Get)
	authed.POST("/cargas/:id/estados", httpx.RequireRole("operador", "administrador"), d.Cargas.CambiarEstado)
	authed.GET("/cargas/:id/trazabilidad", d.Traz.List)
	authed.GET("/cargas/:id/incidencias", d.Incid.List)
	authed.POST("/cargas/:id/incidencias", httpx.RequireRole("operador", "administrador"), d.Incid.Create)
	authed.POST("/incidencias/:id/cerrar", httpx.RequireRole("operador", "administrador"), d.Incid.Cerrar)
	authed.GET("/cargas/:id/documentos", d.Docs.List)
	authed.POST("/cargas/:id/documentos", httpx.RequireRole("operador", "administrador"), d.Docs.Upload)
	authed.GET("/documentos/:id/download", d.Docs.Download)
	authed.GET("/metricas", httpx.RequireRole("administrador", "gerente"), d.Metricas.Get)
	return r
}
