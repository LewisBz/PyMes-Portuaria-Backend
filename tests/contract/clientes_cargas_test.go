package contract_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientesYCargasContratos(t *testing.T) {
	f := setup(t)
	adminTok := loginAdmin(t, f)

	w := f.do(http.MethodPost, "/api/v1/clientes", adminTok, map[string]string{"razon_social": "Naviera Interna"})
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())
	var cli struct {
		ID string `json:"id"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &cli))

	w = f.do(http.MethodGet, "/api/v1/clientes", adminTok, nil)
	require.Equal(t, http.StatusOK, w.Code)

	w = f.do(http.MethodPost, "/api/v1/cargas", adminTok, map[string]string{
		"cliente_id": cli.ID, "referencia": "REF-001",
	})
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())
	var carga struct {
		ID     string `json:"id"`
		Estado string `json:"estado"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &carga))
	require.Equal(t, "registrada", carga.Estado)

	w = f.do(http.MethodGet, "/api/v1/cargas", adminTok, nil)
	require.Equal(t, http.StatusOK, w.Code)
	w = f.do(http.MethodGet, "/api/v1/cargas/"+carga.ID, adminTok, nil)
	require.Equal(t, http.StatusOK, w.Code)
}
