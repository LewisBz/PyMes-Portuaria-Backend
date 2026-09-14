package integration_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestDocumentosUnauthorizedCargaIs404(t *testing.T) {
	e, gdb := openEngine(t)
	orgID := uuid.MustParse("00000000-0000-4000-8000-000000000001")
	cli := insertCliente(t, gdb, orgID)
	carga := insertCarga(t, gdb, orgID, cli)
	other := insertCliente(t, gdb, orgID)
	otherUser := insertClienteUser(t, gdb, orgID, other)
	tok := mintCliente(t, otherUser)
	w := getJSON(e, "/api/v1/cargas/"+carga.String()+"/documentos", tok)
	require.Equal(t, http.StatusNotFound, w.Code)
}
