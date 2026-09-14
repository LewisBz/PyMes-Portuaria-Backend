package contract_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEstadosYTrazabilidadIncluyeGerente(t *testing.T) {
	f := setup(t)
	adminTok := loginAdmin(t, f)
	cli := addCliente(f, "Cliente A")
	carga := addCarga(f, cli.ID, "REF-TRZ")
	gerente := addUser(t, f, "ger@demo.local", "Gerente12!", "gerente", nil)
	gerTok := mint(t, f, gerente)

	w := f.do(http.MethodPost, "/api/v1/cargas/"+carga.ID.String()+"/estados", adminTok, map[string]string{
		"estado": "en_transito",
	})
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	w = f.do(http.MethodGet, "/api/v1/cargas/"+carga.ID.String()+"/trazabilidad", gerTok, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	var ev []map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &ev))
	require.NotEmpty(t, ev)

	w = f.do(http.MethodPost, "/api/v1/cargas/"+carga.ID.String()+"/estados", gerTok, map[string]string{
		"estado": "en_puerto",
	})
	errBodyHasCode(t, w, http.StatusForbidden)
}
