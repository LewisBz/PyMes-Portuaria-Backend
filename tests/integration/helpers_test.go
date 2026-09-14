package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pymes-portuaria/backend/internal/app"
	"github.com/pymes-portuaria/backend/internal/auth"
	"github.com/pymes-portuaria/backend/internal/cargas"
	"github.com/pymes-portuaria/backend/internal/clientes"
	"github.com/pymes-portuaria/backend/internal/documentos"
	"github.com/pymes-portuaria/backend/internal/incidencias"
	"github.com/pymes-portuaria/backend/internal/metricas"
	"github.com/pymes-portuaria/backend/internal/organizacion"
	"github.com/pymes-portuaria/backend/internal/platform/db"
	"github.com/pymes-portuaria/backend/internal/platform/httpx"
	"github.com/pymes-portuaria/backend/internal/platform/storage"
	"github.com/pymes-portuaria/backend/internal/trazabilidad"
	"github.com/pymes-portuaria/backend/internal/usuarios"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const secret = "integration-test-secret-32-bytes-ok"

func skipNoDB(t *testing.T) string {
	t.Helper()
	u := os.Getenv("TEST_DATABASE_URL")
	if u == "" {
		u = os.Getenv("DATABASE_URL")
	}
	if u == "" {
		t.Skip("DATABASE_URL / TEST_DATABASE_URL not set")
	}
	return u
}

func migrationsDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "migrations"))
}

func openEngine(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	u := skipNoDB(t)
	require.NoError(t, db.Migrate(u, migrationsDir(t)))
	gdb, err := db.Open(u)
	require.NoError(t, err)
	orgRepo := &organizacion.Postgres{DB: gdb}
	userRepo := &usuarios.Postgres{DB: gdb}
	cliRepo := &clientes.Postgres{DB: gdb}
	cargaRepo := &cargas.Postgres{DB: gdb}
	evRepo := &trazabilidad.Postgres{DB: gdb}
	incRepo := &incidencias.Postgres{DB: gdb}
	docRepo := &documentos.Postgres{DB: gdb}
	cargaSvc := &cargas.Service{Repo: cargaRepo, Clientes: cliRepo, Events: evRepo, DB: gdb}
	gin.SetMode(gin.TestMode)
	e := app.NewEngine(app.Deps{
		JWTSecret: secret,
		Orgs:      orgRepo,
		Users:     userRepo,
		Auth:      &auth.Handler{Svc: &auth.Service{Users: userRepo, Orgs: orgRepo, Secret: secret, Expiry: time.Hour}},
		Usuarios:  &usuarios.Handler{Svc: &usuarios.Service{Repo: userRepo, Cost: 12}},
		Clientes:  &clientes.Handler{Svc: &clientes.Service{Repo: cliRepo}},
		Cargas:    &cargas.Handler{Svc: cargaSvc},
		Traz:      &trazabilidad.Handler{Svc: &trazabilidad.Service{Repo: evRepo, Carga: cargaSvc}},
		Incid:     &incidencias.Handler{Svc: &incidencias.Service{Repo: incRepo, Cargas: cargaSvc, Events: evRepo, DB: gdb}},
		Docs:      &documentos.Handler{Svc: &documentos.Service{Repo: docRepo, Store: &storage.Local{Root: t.TempDir()}, Cargas: cargaSvc, Events: evRepo}},
		Metricas:  &metricas.Handler{Svc: &metricas.Service{DB: gdb}},
	})
	return e, gdb
}

func postJSON(e *gin.Engine, path, token string, body any) *httptest.ResponseRecorder {
	var r *bytes.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		r = bytes.NewReader(b)
	} else {
		r = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(http.MethodPost, path, r)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	return w
}

func getJSON(e *gin.Engine, path, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	return w
}

func loginAdmin(t *testing.T, e *gin.Engine) string {
	t.Helper()
	w := postJSON(e, "/api/v1/auth/login", "", map[string]string{"email": "admin@demo.local", "password": "DemoAdmin12!"})
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	var res struct {
		AccessToken string `json:"access_token"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	return res.AccessToken
}

func mintCliente(t *testing.T, u usuarios.Usuario) string {
	t.Helper()
	claims := httpx.TokenClaims{
		UserID: u.ID.String(), OrganizacionID: u.OrganizacionID.String(), Rol: string(u.Rol),
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour))},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString([]byte(secret))
	require.NoError(t, err)
	return s
}

func insertCliente(t *testing.T, gdb *gorm.DB, orgID uuid.UUID) uuid.UUID {
	t.Helper()
	id := uuid.New()
	require.NoError(t, gdb.Exec(
		`INSERT INTO clientes (id, organizacion_id, razon_social, activo) VALUES (?, ?, ?, true)`,
		id, orgID, "Cli-"+id.String()[:8],
	).Error)
	return id
}

func insertCarga(t *testing.T, gdb *gorm.DB, orgID, clienteID uuid.UUID) uuid.UUID {
	t.Helper()
	id := uuid.New()
	ref := "REF-" + id.String()[:8]
	require.NoError(t, gdb.Exec(
		`INSERT INTO cargas (id, organizacion_id, cliente_id, referencia, estado) VALUES (?, ?, ?, ?, 'registrada')`,
		id, orgID, clienteID, ref,
	).Error)
	return id
}

func insertClienteUser(t *testing.T, gdb *gorm.DB, orgID, clienteID uuid.UUID) usuarios.Usuario {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte("Cliente12!"), 12)
	require.NoError(t, err)
	id := uuid.New()
	email := "c-" + id.String()[:8] + "@demo.local"
	require.NoError(t, gdb.Exec(
		`INSERT INTO usuarios (id, organizacion_id, email, password_hash, nombre, rol, cliente_id, activo) VALUES (?, ?, ?, ?, ?, 'cliente', ?, true)`,
		id, orgID, email, string(hash), "Cliente", clienteID,
	).Error)
	return usuarios.Usuario{ID: id, OrganizacionID: orgID, Email: email, Rol: usuarios.RolCliente, ClienteID: &clienteID, Activo: true}
}
