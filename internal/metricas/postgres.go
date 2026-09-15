package metricas

// SQL aggregates for GET /metricas live in service.go so handlers never import GORM.
// Queries (all scoped by organizacion_id):
//   cargas_activas: estado NOT IN (entregada, cancelada)
//   cargas_retrasadas: fecha_comprometida < NOW() AND estado <> entregada
//   incidencias_abiertas: estado = abierta
//   tiempo_promedio_total_entrega_horas: AVG(fecha_entrega - created_at) hours, entregada only
