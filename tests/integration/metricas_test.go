package integration_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetricasAggregation(t *testing.T) {
	e, _ := openEngine(t)
	adminTok := loginAdmin(t, e)
	w := getJSON(e, "/api/v1/metricas", adminTok)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	require.Contains(t, w.Body.String(), "cargas_activas")
	require.Contains(t, w.Body.String(), "tiempo_promedio_total_entrega_horas")
}
