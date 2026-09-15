package integration_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthRolesLogin(t *testing.T) {
	e, _ := openEngine(t)
	_ = loginAdmin(t, e)
	w := postJSON(e, "/api/v1/auth/login", "", map[string]string{"email": "admin@demo.local", "password": "bad"})
	require.Equal(t, http.StatusUnauthorized, w.Code)
}
