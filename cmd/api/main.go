package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pymes-portuaria/backend/internal/app"
	"github.com/pymes-portuaria/backend/internal/auth"
	"github.com/pymes-portuaria/backend/internal/cargas"
	"github.com/pymes-portuaria/backend/internal/clientes"
	"github.com/pymes-portuaria/backend/internal/config"
	"github.com/pymes-portuaria/backend/internal/documentos"
	"github.com/pymes-portuaria/backend/internal/incidencias"
	"github.com/pymes-portuaria/backend/internal/metricas"
	"github.com/pymes-portuaria/backend/internal/organizacion"
	"github.com/pymes-portuaria/backend/internal/platform/db"
	"github.com/pymes-portuaria/backend/internal/platform/storage"
	"github.com/pymes-portuaria/backend/internal/trazabilidad"
	"github.com/pymes-portuaria/backend/internal/usuarios"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Migrate(cfg.DatabaseURL, cfg.Migrations); err != nil {
		log.Fatal(err)
	}
	gdb, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(cfg.UploadDir, 0o755); err != nil {
		log.Fatal(err)
	}

	orgRepo := &organizacion.Postgres{DB: gdb}
	userRepo := &usuarios.Postgres{DB: gdb}
	cliRepo := &clientes.Postgres{DB: gdb}
	cargaRepo := &cargas.Postgres{DB: gdb}
	evRepo := &trazabilidad.Postgres{DB: gdb}
	incRepo := &incidencias.Postgres{DB: gdb}
	docRepo := &documentos.Postgres{DB: gdb}

	userSvc := &usuarios.Service{Repo: userRepo, Cost: cfg.BcryptCost}
	authSvc := &auth.Service{Users: userRepo, Orgs: orgRepo, Secret: cfg.JWTSecret, Expiry: cfg.JWTExpiry}
	cliSvc := &clientes.Service{Repo: cliRepo}
	cargaSvc := &cargas.Service{Repo: cargaRepo, Clientes: cliRepo, Events: evRepo, DB: gdb}
	trazSvc := &trazabilidad.Service{Repo: evRepo, Carga: cargaSvc}
	incSvc := &incidencias.Service{Repo: incRepo, Cargas: cargaSvc, Events: evRepo, DB: gdb}
	docSvc := &documentos.Service{
		Repo: docRepo, Store: &storage.Local{Root: cfg.UploadDir},
		Cargas: cargaSvc, Events: evRepo,
	}
	metSvc := &metricas.Service{DB: gdb}

	engine := app.NewEngine(app.Deps{
		JWTSecret: cfg.JWTSecret,
		Orgs:      orgRepo,
		Users:     userRepo,
		Auth:      &auth.Handler{Svc: authSvc},
		Usuarios:  &usuarios.Handler{Svc: userSvc},
		Clientes:  &clientes.Handler{Svc: cliSvc},
		Cargas:    &cargas.Handler{Svc: cargaSvc},
		Traz:      &trazabilidad.Handler{Svc: trazSvc},
		Incid:     &incidencias.Handler{Svc: incSvc},
		Docs:      &documentos.Handler{Svc: docSvc},
		Metricas:  &metricas.Handler{Svc: metSvc},
	})

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, engine); err != nil {
		log.Fatal(err)
	}
}
