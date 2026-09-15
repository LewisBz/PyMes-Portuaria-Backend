package integration_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestIncidenciasOpenClose(t *testing.T) {
	e, gdb := openEngine(t)
	adminTok := loginAdmin(t, e)
	orgID := uuid.MustParse("00000000-0000-4000-8000-000000000001")
	cli := insertCliente(t, gdb, orgID)
	carga := insertCarga(t, gdb, orgID, cli)
	w := postJSON(e, "/api/v1/cargas/"+carga.String()+"/incidencias", adminTok, map[string]string{
		"titulo": "t", "descripcion": "d", "tipo": "otro",
	})
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())
	var inc struct {
		ID string `json:"id"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &inc))
	w = postJSON(e, "/api/v1/incidencias/"+inc.ID+"/cerrar", adminTok, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
}
