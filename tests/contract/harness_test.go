package contract_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
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
	"github.com/pymes-portuaria/backend/internal/platform/httpx"
	"github.com/pymes-portuaria/backend/internal/platform/storage"
	"github.com/pymes-portuaria/backend/internal/trazabilidad"
	"github.com/pymes-portuaria/backend/internal/usuarios"
	"github.com/pymes-portuaria/backend/pkg/httperr"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

const testSecret = "contract-test-secret-32-bytes-min"

type mem struct {
	mu     sync.Mutex
	org    organizacion.Organizacion
	users  map[uuid.UUID]*usuarios.Usuario
	cli    map[uuid.UUID]*clientes.Cliente
	cargas map[uuid.UUID]*cargas.Carga
	ev     []trazabilidad.Evento
	inc    map[uuid.UUID]*incidencias.Incidencia
	docs   map[uuid.UUID]*documentos.Documento
}

func newMem() *mem {
	orgID := uuid.MustParse("00000000-0000-4000-8000-000000000001")
	hash, _ := bcrypt.GenerateFromPassword([]byte("DemoAdmin12!"), 12)
	adminID := uuid.MustParse("00000000-0000-4000-8000-000000000002")
	m := &mem{
		org:    organizacion.Organizacion{ID: orgID, Nombre: "Demo PyME", Activa: true},
		users:  map[uuid.UUID]*usuarios.Usuario{},
		cli:    map[uuid.UUID]*clientes.Cliente{},
		cargas: map[uuid.UUID]*cargas.Carga{},
		inc:    map[uuid.UUID]*incidencias.Incidencia{},
		docs:   map[uuid.UUID]*documentos.Documento{},
	}
	m.users[adminID] = &usuarios.Usuario{
		ID: adminID, OrganizacionID: orgID, Email: "admin@demo.local",
		PasswordHash: string(hash), Nombre: "Admin", Rol: usuarios.RolAdministrador, Activo: true,
	}
	return m
}

func (m *mem) GetByID(_ context.Context, id uuid.UUID) (*organizacion.Organizacion, error) {
	if m.org.ID != id {
		return nil, httperr.NotFound("organización no encontrada")
	}
	o := m.org
	return &o, nil
}
func (m *mem) IsActiva(_ context.Context, id uuid.UUID) (bool, error) {
	o, err := m.GetByID(context.Background(), id)
	if err != nil {
		return false, err
	}
	return o.Activa, nil
}

func (m *mem) userGetByID(_ context.Context, orgID, id uuid.UUID) (*usuarios.Usuario, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.users[id]
	if !ok || u.OrganizacionID != orgID {
		return nil, httperr.NotFound("usuario no encontrado")
	}
	cp := *u
	return &cp, nil
}
func (m *mem) GetByEmail(_ context.Context, orgID uuid.UUID, email string) (*usuarios.Usuario, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, u := range m.users {
		if u.OrganizacionID == orgID && u.Email == email {
			cp := *u
			return &cp, nil
		}
	}
	return nil, httperr.NotFound("usuario no encontrado")
}
func (m *mem) GetByEmailAnyOrg(_ context.Context, email string) (*usuarios.Usuario, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, u := range m.users {
		if u.Email == email {
			cp := *u
			return &cp, nil
		}
	}
	return nil, httperr.Unauthorized("credenciales inválidas")
}
func (m *mem) List(_ context.Context, orgID uuid.UUID) ([]usuarios.Usuario, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []usuarios.Usuario
	for _, u := range m.users {
		if u.OrganizacionID == orgID {
			out = append(out, *u)
		}
	}
	return out, nil
}
func (m *mem) Create(_ context.Context, u *usuarios.Usuario) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *u
	m.users[u.ID] = &cp
	return nil
}
func (m *mem) Update(_ context.Context, u *usuarios.Usuario) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *u
	m.users[u.ID] = &cp
	return nil
}

func (m *mem) cliCreate(_ context.Context, c *clientes.Cliente) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *c
	m.cli[c.ID] = &cp
	return nil
}
func (m *mem) cliList(_ context.Context, orgID uuid.UUID) ([]clientes.Cliente, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []clientes.Cliente
	for _, c := range m.cli {
		if c.OrganizacionID == orgID {
			out = append(out, *c)
		}
	}
	return out, nil
}
func (m *mem) cliGet(_ context.Context, orgID, id uuid.UUID) (*clientes.Cliente, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.cli[id]
	if !ok || c.OrganizacionID != orgID {
		return nil, httperr.NotFound("cliente no encontrado")
	}
	cp := *c
	return &cp, nil
}

func (m *mem) cargaCreate(_ context.Context, c *cargas.Carga) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *c
	m.cargas[c.ID] = &cp
	return nil
}
func (m *mem) cargaList(_ context.Context, orgID uuid.UUID, clienteID *uuid.UUID, estado *cargas.Estado) ([]cargas.Carga, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []cargas.Carga
	for _, c := range m.cargas {
		if c.OrganizacionID != orgID {
			continue
		}
		if clienteID != nil && c.ClienteID != *clienteID {
			continue
		}
		if estado != nil && c.Estado != *estado {
			continue
		}
		out = append(out, *c)
	}
	return out, nil
}
func (m *mem) cargaGet(_ context.Context, orgID, id uuid.UUID) (*cargas.Carga, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.cargas[id]
	if !ok || c.OrganizacionID != orgID {
		return nil, httperr.NotFound("carga no encontrada")
	}
	cp := *c
	return &cp, nil
}
func (m *mem) cargaUpdate(_ context.Context, c *cargas.Carga) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *c
	m.cargas[c.ID] = &cp
	return nil
}

func (m *mem) Append(_ context.Context, e trazabilidad.Evento) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ev = append(m.ev, e)
	return nil
}
func (m *mem) ListByCarga(_ context.Context, orgID, cargaID uuid.UUID) ([]trazabilidad.Evento, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []trazabilidad.Evento
	for _, e := range m.ev {
		if e.OrganizacionID == orgID && e.CargaID == cargaID {
			out = append(out, e)
		}
	}
	return out, nil
}

func (m *mem) incCreate(_ context.Context, i *incidencias.Incidencia) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *i
	m.inc[i.ID] = &cp
	return nil
}
func (m *mem) incList(_ context.Context, orgID, cargaID uuid.UUID) ([]incidencias.Incidencia, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []incidencias.Incidencia
	for _, i := range m.inc {
		if i.OrganizacionID == orgID && i.CargaID == cargaID {
			out = append(out, *i)
		}
	}
	return out, nil
}
func (m *mem) incGet(_ context.Context, orgID, id uuid.UUID) (*incidencias.Incidencia, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	i, ok := m.inc[id]
	if !ok || i.OrganizacionID != orgID {
		return nil, httperr.NotFound("incidencia no encontrada")
	}
	cp := *i
	return &cp, nil
}
func (m *mem) incUpdate(_ context.Context, i *incidencias.Incidencia) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *i
	m.inc[i.ID] = &cp
	return nil
}

func (m *mem) docCreate(_ context.Context, d *documentos.Documento) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *d
	m.docs[d.ID] = &cp
	return nil
}
func (m *mem) docList(_ context.Context, orgID, cargaID uuid.UUID, onlyVisible bool) ([]documentos.Documento, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []documentos.Documento
	for _, d := range m.docs {
		if d.OrganizacionID != orgID || d.CargaID != cargaID {
			continue
		}
		if onlyVisible && !d.VisibleCliente {
			continue
		}
		out = append(out, *d)
	}
	return out, nil
}
func (m *mem) docGet(_ context.Context, orgID, id uuid.UUID) (*documentos.Documento, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	d, ok := m.docs[id]
	if !ok || d.OrganizacionID != orgID {
		return nil, httperr.NotFound("documento no encontrado")
	}
	cp := *d
	return &cp, nil
}

type userAdapter struct{ m *mem }
type cliAdapter struct{ m *mem }
type cargaAdapter struct{ m *mem }
type incAdapter struct{ m *mem }
type docAdapter struct{ m *mem }

func (a userAdapter) GetByID(ctx context.Context, orgID, id uuid.UUID) (*usuarios.Usuario, error) {
	return a.m.userGetByID(ctx, orgID, id)
}
func (a userAdapter) GetByEmail(ctx context.Context, orgID uuid.UUID, email string) (*usuarios.Usuario, error) {
	return a.m.GetByEmail(ctx, orgID, email)
}
func (a userAdapter) GetByEmailAnyOrg(ctx context.Context, email string) (*usuarios.Usuario, error) {
	return a.m.GetByEmailAnyOrg(ctx, email)
}
func (a userAdapter) List(ctx context.Context, orgID uuid.UUID) ([]usuarios.Usuario, error) {
	return a.m.List(ctx, orgID)
}
func (a userAdapter) Create(ctx context.Context, u *usuarios.Usuario) error {
	return a.m.Create(ctx, u)
}
func (a userAdapter) Update(ctx context.Context, u *usuarios.Usuario) error {
	return a.m.Update(ctx, u)
}

func (a cliAdapter) Create(ctx context.Context, c *clientes.Cliente) error {
	return a.m.cliCreate(ctx, c)
}
func (a cliAdapter) List(ctx context.Context, orgID uuid.UUID) ([]clientes.Cliente, error) {
	return a.m.cliList(ctx, orgID)
}
func (a cliAdapter) GetByID(ctx context.Context, orgID, id uuid.UUID) (*clientes.Cliente, error) {
	return a.m.cliGet(ctx, orgID, id)
}

func (a cargaAdapter) Create(ctx context.Context, c *cargas.Carga) error {
	return a.m.cargaCreate(ctx, c)
}
func (a cargaAdapter) List(ctx context.Context, orgID uuid.UUID, clienteID *uuid.UUID, estado *cargas.Estado) ([]cargas.Carga, error) {
	return a.m.cargaList(ctx, orgID, clienteID, estado)
}
func (a cargaAdapter) GetByID(ctx context.Context, orgID, id uuid.UUID) (*cargas.Carga, error) {
	return a.m.cargaGet(ctx, orgID, id)
}
func (a cargaAdapter) Update(ctx context.Context, c *cargas.Carga) error {
	return a.m.cargaUpdate(ctx, c)
}

func (a incAdapter) Create(ctx context.Context, i *incidencias.Incidencia) error {
	return a.m.incCreate(ctx, i)
}
func (a incAdapter) ListByCarga(ctx context.Context, orgID, cargaID uuid.UUID) ([]incidencias.Incidencia, error) {
	return a.m.incList(ctx, orgID, cargaID)
}
func (a incAdapter) GetByID(ctx context.Context, orgID, id uuid.UUID) (*incidencias.Incidencia, error) {
	return a.m.incGet(ctx, orgID, id)
}
func (a incAdapter) Update(ctx context.Context, i *incidencias.Incidencia) error {
	return a.m.incUpdate(ctx, i)
}

func (a docAdapter) Create(ctx context.Context, d *documentos.Documento) error {
	return a.m.docCreate(ctx, d)
}
func (a docAdapter) ListByCarga(ctx context.Context, orgID, cargaID uuid.UUID, onlyVisible bool) ([]documentos.Documento, error) {
	return a.m.docList(ctx, orgID, cargaID, onlyVisible)
}
func (a docAdapter) GetByID(ctx context.Context, orgID, id uuid.UUID) (*documentos.Documento, error) {
	return a.m.docGet(ctx, orgID, id)
}

type fixture struct {
	engine *gin.Engine
	mem    *mem
	dir    string
}

func setup(t *testing.T) *fixture {
	t.Helper()
	gin.SetMode(gin.TestMode)
	m := newMem()
	users := userAdapter{m}
	clis := cliAdapter{m}
	cargaR := cargaAdapter{m}
	cargaSvc := &cargas.Service{Repo: cargaR, Clientes: clis, Events: m, DB: nil}
	incSvc := &incidencias.Service{Repo: incAdapter{m}, Cargas: cargaSvc, Events: m, DB: nil}
	dir := t.TempDir()
	docSvc := &documentos.Service{
		Repo: docAdapter{m}, Store: &storage.Local{Root: dir}, Cargas: cargaSvc, Events: m,
	}
	eng := app.NewEngine(app.Deps{
		JWTSecret: testSecret,
		Orgs:      m,
		Users:     users,
		Auth:      &auth.Handler{Svc: &auth.Service{Users: users, Orgs: m, Secret: testSecret, Expiry: time.Hour}},
		Usuarios:  &usuarios.Handler{Svc: &usuarios.Service{Repo: users, Cost: 12}},
		Clientes:  &clientes.Handler{Svc: &clientes.Service{Repo: clis}},
		Cargas:    &cargas.Handler{Svc: cargaSvc},
		Traz:      &trazabilidad.Handler{Svc: &trazabilidad.Service{Repo: m, Carga: cargaSvc}},
		Incid:     &incidencias.Handler{Svc: incSvc},
		Docs:      &documentos.Handler{Svc: docSvc},
		Metricas:  &metricas.Handler{Svc: stubMetrics{}},
	})
	return &fixture{engine: eng, mem: m, dir: dir}
}

type stubMetrics struct{}

func (stubMetrics) Get(context.Context, uuid.UUID) (metricas.Snapshot, error) {
	return metricas.Snapshot{}, nil
}

func (f *fixture) do(method, path, token string, body any) *httptest.ResponseRecorder {
	var rdr io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		rdr = bytes.NewReader(b)
	}
	req := httptest.NewRequest(method, path, rdr)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	f.engine.ServeHTTP(w, req)
	return w
}

func decodeErr(t *testing.T, w *httptest.ResponseRecorder) httperr.Error {
	t.Helper()
	var e httperr.Error
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &e))
	require.Equal(t, w.Code, e.Code)
	require.NotEmpty(t, e.Message)
	return e
}

func loginAdmin(t *testing.T, f *fixture) string {
	t.Helper()
	w := f.do(http.MethodPost, "/api/v1/auth/login", "", map[string]string{
		"email": "admin@demo.local", "password": "DemoAdmin12!",
	})
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	var res struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	require.Equal(t, "Bearer", res.TokenType)
	return res.AccessToken
}

func claimsOf(t *testing.T, token string) httpx.TokenClaims {
	t.Helper()
	tc := &httpx.TokenClaims{}
	_, err := jwt.ParseWithClaims(token, tc, func(tok *jwt.Token) (any, error) {
		return []byte(testSecret), nil
	})
	require.NoError(t, err)
	require.NotEmpty(t, tc.UserID)
	require.NotEmpty(t, tc.OrganizacionID)
	require.NotEmpty(t, tc.Rol)
	return *tc
}

func mint(t *testing.T, f *fixture, u *usuarios.Usuario) string {
	t.Helper()
	claims := httpx.TokenClaims{
		UserID: u.ID.String(), OrganizacionID: u.OrganizacionID.String(), Rol: string(u.Rol),
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour))},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString([]byte(testSecret))
	require.NoError(t, err)
	return s
}

func addUser(t *testing.T, f *fixture, email, pass, rol string, clienteID *uuid.UUID) *usuarios.Usuario {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), 12)
	require.NoError(t, err)
	u := &usuarios.Usuario{
		ID: uuid.New(), OrganizacionID: f.mem.org.ID, Email: email,
		PasswordHash: string(hash), Nombre: email, Rol: usuarios.Rol(rol), ClienteID: clienteID, Activo: true,
	}
	require.NoError(t, f.mem.Create(context.Background(), u))
	return u
}

func addCliente(f *fixture, nombre string) *clientes.Cliente {
	c := &clientes.Cliente{ID: uuid.New(), OrganizacionID: f.mem.org.ID, RazonSocial: nombre, Activo: true}
	_ = f.mem.cliCreate(context.Background(), c)
	return c
}

func addCarga(f *fixture, clienteID uuid.UUID, ref string) *cargas.Carga {
	c := &cargas.Carga{
		ID: uuid.New(), OrganizacionID: f.mem.org.ID, ClienteID: clienteID,
		Referencia: ref, Estado: cargas.EstadoRegistrada,
	}
	_ = f.mem.cargaCreate(context.Background(), c)
	return c
}

func multipartPDF(t *testing.T, visible bool) (*bytes.Buffer, string) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", "doc.pdf")
	require.NoError(t, err)
	_, _ = fw.Write([]byte("%PDF-1.4 test"))
	_ = w.WriteField("visible_cliente", map[bool]string{true: "true", false: "false"}[visible])
	require.NoError(t, w.Close())
	return &buf, w.FormDataContentType()
}

func errBodyHasCode(t *testing.T, w *httptest.ResponseRecorder, code int) {
	t.Helper()
	require.Equal(t, code, w.Code, w.Body.String())
	e := decodeErr(t, w)
	require.Equal(t, code, e.Code)
}

func TestLoginEmptyBody400(t *testing.T) {
	f := setup(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(""))
	w := httptest.NewRecorder()
	f.engine.ServeHTTP(w, req)
	errBodyHasCode(t, w, http.StatusBadRequest)
}
