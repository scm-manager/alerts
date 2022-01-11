package api

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCors(t *testing.T) {
	cors := Cors(CreateOkEndpoint())
	r, err := http.NewRequest("GET", "/something", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	cors.ServeHTTP(w, r)

	assert.Equal(t, "*", w.Result().Header.Get("Access-Control-Allow-Origin"))
}
