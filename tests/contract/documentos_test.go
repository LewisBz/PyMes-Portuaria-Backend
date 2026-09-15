package contract_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDocumentosContrato(t *testing.T) {
	f := setup(t)
	adminTok := loginAdmin(t, f)
	cli := addCliente(f, "Cliente Doc")
	carga := addCarga(f, cli.ID, "REF-DOC")
	clienteUser := addUser(t, f, "doc-cli@demo.local", "Cliente12!", "cliente", &cli.ID)
	cliTok := mint(t, f, clienteUser)

	body, ctype := multipartPDF(t, true)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cargas/"+carga.ID.String()+"/documentos", body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("Authorization", "Bearer "+adminTok)
	w := httptest.NewRecorder()
	f.engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())
	var doc struct {
		ID string `json:"id"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &doc))

	list := f.do(http.MethodGet, "/api/v1/cargas/"+carga.ID.String()+"/documentos", cliTok, nil)
	require.Equal(t, http.StatusOK, list.Code)

	dl := f.do(http.MethodGet, "/api/v1/documentos/"+doc.ID+"/download", cliTok, nil)
	require.Equal(t, http.StatusOK, dl.Code)
	require.Contains(t, dl.Body.String(), "%PDF")
}
