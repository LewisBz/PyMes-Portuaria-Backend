package contract_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIncidenciasContrato(t *testing.T) {
	f := setup(t)
	adminTok := loginAdmin(t, f)
	cli := addCliente(f, "Cliente Inc")
	carga := addCarga(f, cli.ID, "REF-INC")
	clienteUser := addUser(t, f, "inc-cli@demo.local", "Cliente12!", "cliente", &cli.ID)
	cliTok := mint(t, f, clienteUser)

	w := f.do(http.MethodPost, "/api/v1/cargas/"+carga.ID.String()+"/incidencias", adminTok, map[string]string{
		"titulo": "Demora", "descripcion": "Atraso en patio", "tipo": "demora",
	})
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())
	var inc struct {
		ID     string `json:"id"`
		Estado string `json:"estado"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &inc))
	require.Equal(t, "abierta", inc.Estado)

	w = f.do(http.MethodGet, "/api/v1/cargas/"+carga.ID.String()+"/incidencias", adminTok, nil)
	require.Equal(t, http.StatusOK, w.Code)

	w = f.do(http.MethodPost, "/api/v1/incidencias/"+inc.ID+"/cerrar", adminTok, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	w = f.do(http.MethodPost, "/api/v1/cargas/"+carga.ID.String()+"/incidencias", cliTok, map[string]string{
		"titulo": "No", "descripcion": "No",
	})
	errBodyHasCode(t, w, http.StatusForbidden)
}
