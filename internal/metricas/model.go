package metricas

type Snapshot struct {
	CargasActivas                   int      `json:"cargas_activas"`
	CargasRetrasadas                int      `json:"cargas_retrasadas"`
	IncidenciasAbiertas             int      `json:"incidencias_abiertas"`
	TiempoPromedioTotalEntregaHoras *float64 `json:"tiempo_promedio_total_entrega_horas"`
}
