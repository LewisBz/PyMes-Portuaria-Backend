package contract_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoginAndMeClaims(t *testing.T) {
	f := setup(t)
	token := loginAdmin(t, f)
	cl := claimsOf(t, token)
	require.Equal(t, "administrador", cl.Rol)
	require.Equal(t, "00000000-0000-4000-8000-000000000002", cl.UserID)
	require.Equal(t, "00000000-0000-4000-8000-000000000001", cl.OrganizacionID)

	w := f.do(http.MethodGet, "http://localhost/api/v1/auth/me", token, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	var u map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &u))
	require.Equal(t, "admin@demo.local", u["email"])
	_, hasHash := u["password_hash"]
	require.False(t, hasHash)
}

func TestLoginBadPassword401(t *testing.T) {
	f := setup(t)
	w := f.do(http.MethodPost, "/api/v1/auth/login", "", map[string]string{
		"email": "admin@demo.local", "password": "wrong-password",
	})
	errBodyHasCode(t, w, http.StatusUnauthorized)
}
