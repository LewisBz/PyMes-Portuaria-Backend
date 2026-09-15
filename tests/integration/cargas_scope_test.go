package integration_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCargasScope404Never403(t *testing.T) {
	e, gdb := openEngine(t)
	orgID := uuid.MustParse("00000000-0000-4000-8000-000000000001")
	a := insertCliente(t, gdb, orgID)
	b := insertCliente(t, gdb, orgID)
	carga := insertCarga(t, gdb, orgID, a)
	userB := insertClienteUser(t, gdb, orgID, b)
	tok := mintCliente(t, userB)
	w := getJSON(e, "/api/v1/cargas/"+carga.String(), tok)
	require.Equal(t, http.StatusNotFound, w.Code, w.Body.String())
	require.Contains(t, w.Body.String(), `"code":404`)
	w = postJSON(e, "/api/v1/cargas", tok, map[string]string{"cliente_id": a.String(), "referencia": "X"})
	require.Equal(t, http.StatusForbidden, w.Code)
}
