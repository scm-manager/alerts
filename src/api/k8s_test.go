package api

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateOkEndpoint(t *testing.T) {
	handler := CreateOkEndpoint()

	r, err := http.NewRequest("GET", "/ready", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}
