package contract_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUsuariosAdminVsCliente(t *testing.T) {
	f := setup(t)
	adminTok := loginAdmin(t, f)

	w := f.do(http.MethodGet, "/api/v1/usuarios", adminTok, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	w = f.do(http.MethodPost, "/api/v1/usuarios", adminTok, map[string]any{
		"email": "op@demo.local", "password": "Operador12!", "nombre": "Op", "rol": "operador",
	})
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())

	cli := addCliente(f, "ACME")
	clienteUser := addUser(t, f, "cli@demo.local", "Cliente12!", "cliente", &cli.ID)
	cliTok := mint(t, f, clienteUser)

	w = f.do(http.MethodGet, "/api/v1/usuarios", cliTok, nil)
	errBodyHasCode(t, w, http.StatusForbidden)

	w = f.do(http.MethodPost, "/api/v1/usuarios", cliTok, map[string]any{
		"email": "x@demo.local", "password": "Password12!", "nombre": "X", "rol": "operador",
	})
	errBodyHasCode(t, w, http.StatusForbidden)

	adminID := uuid.MustParse("00000000-0000-4000-8000-000000000002")
	w = f.do(http.MethodPatch, "/api/v1/usuarios/"+adminID.String(), cliTok, map[string]any{"activo": false})
	errBodyHasCode(t, w, http.StatusForbidden)

	w = f.do(http.MethodPatch, "/api/v1/usuarios/"+adminID.String(), adminTok, map[string]any{"nombre": "Admin 2"})
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
}
