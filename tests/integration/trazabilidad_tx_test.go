package integration_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestTrazabilidadOneEventPerStateChange(t *testing.T) {
	e, gdb := openEngine(t)
	orgID := uuid.MustParse("00000000-0000-4000-8000-000000000001")
	cli := insertCliente(t, gdb, orgID)
	carga := insertCarga(t, gdb, orgID, cli)
	adminTok := loginAdmin(t, e)
	w := postJSON(e, "/api/v1/cargas/"+carga.String()+"/estados", adminTok, map[string]string{"estado": "en_transito"})
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	w = postJSON(e, "/api/v1/cargas/"+carga.String()+"/estados", adminTok, map[string]string{"estado": "en_transito"})
	require.Equal(t, http.StatusConflict, w.Code)
	w = getJSON(e, "/api/v1/cargas/"+carga.String()+"/trazabilidad", adminTok)
	require.Equal(t, http.StatusOK, w.Code)
	var ev []map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &ev))
	n := 0
	for _, item := range ev {
		if item["tipo"] == "estado_cambiado" {
			n++
		}
	}
	require.Equal(t, 1, n)
}
