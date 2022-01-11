package api

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInstrumentHandler(t *testing.T) {
	handler := InstrumentHandler(CreateOkEndpoint())

	r, err := http.NewRequest("GET", "/ok", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	metricsHandler := promhttp.Handler()
	metricsHandler.ServeHTTP(w, r)

	defer w.Result().Body.Close()
	data, err := ioutil.ReadAll(w.Result().Body)
	assert.NoError(t, err)

	metrics := string(data)
	assert.Contains(t, metrics, "api_requests_total{code=\"200\",method=\"get\"} 1")
	assert.Contains(t, metrics, "in_flight_requests 0")
	assert.Contains(t, metrics, "push_request_size_bytes_bucket{le=\"200\"} 1")
	assert.Contains(t, metrics, "response_duration_seconds_bucket{method=\"get\",le=\"10\"} 1")
}
