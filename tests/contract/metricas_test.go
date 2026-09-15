package contract_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetricasGerente200Cliente403(t *testing.T) {
	f := setup(t)
	gerente := addUser(t, f, "mger@demo.local", "Gerente12!", "gerente", nil)
	gerTok := mint(t, f, gerente)
	w := f.do(http.MethodGet, "/api/v1/metricas", gerTok, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	cli := addCliente(f, "C")
	cu := addUser(t, f, "mcli@demo.local", "Cliente12!", "cliente", &cli.ID)
	cliTok := mint(t, f, cu)
	w = f.do(http.MethodGet, "/api/v1/metricas", cliTok, nil)
	errBodyHasCode(t, w, http.StatusForbidden)
}
