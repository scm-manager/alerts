package api

import (
	"bytes"
	"encoding/json"
	"github.com/scm-manager/alerts/src/alert"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAlertsEndpoint_ServeHTTP(t *testing.T) {
	alerts := make(alert.Alerts)
	alerts["core"] = []alert.Alert{{
		Title:            "Defect",
		Description:      "Incredible Defect",
		IssuedAt:         alert.IssuedAt{Time: time.Now()},
		AffectedVersions: alert.MustParseVersionRange("<2.27.3"),
	}}

	endpoint := CreateAlertsEndpoint(alerts)

	body := alertsRequest{
		Version: "2.27.1",
		Os:      "Linux",
		Arch:    "amd64",
		Java:    "1.8.0_121",
	}

	data, err := json.Marshal(&body)
	assert.NoError(t, err)

	request, err := http.NewRequest("POST", "/api/v1/alerts", bytes.NewReader(data))
	assert.NoError(t, err)
	request.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	endpoint.ServeHTTP(w, request)

	assert.Equal(t, 200, w.Result().StatusCode)

	defer w.Result().Body.Close()
	data, err = ioutil.ReadAll(w.Result().Body)
	assert.NoError(t, err)

	responseBody := alertsResponse{}
	err = json.Unmarshal(data, &responseBody)
	assert.NoError(t, err)

	assert.Len(t, responseBody.Alerts, 1)
	assert.Equal(t, responseBody.Alerts[0].Title, "Defect")
}
