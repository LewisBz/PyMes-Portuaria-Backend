package cargas

import (
	"testing"
	"time"
)

func TestRetrasadaDerived(t *testing.T) {
	past := time.Now().Add(-24 * time.Hour)
	c := Carga{Estado: EstadoEnTransito, FechaComprometida: &past}
	if !c.ToDTO().Retrasada {
		t.Fatal("expected retrasada")
	}
	c.Estado = EstadoEntregada
	if c.ToDTO().Retrasada {
		t.Fatal("entregada must not be retrasada")
	}
}
