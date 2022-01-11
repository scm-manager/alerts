package api

import (
	"bytes"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/scm-manager/alerts/src/alert"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func requestAlerts(t *testing.T, alerts alert.Alerts, body interface{}) *http.Response {
	endpoint := CreateAlertsEndpoint(alerts)

	data, err := json.Marshal(&body)
	assert.NoError(t, err)

	request, err := http.NewRequest("POST", "/api/v1/alerts", bytes.NewReader(data))
	assert.NoError(t, err)
	request.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	endpoint.ServeHTTP(w, request)

	return w.Result()
}

func callAlertsEndpoint(t *testing.T, alerts alert.Alerts, body alertsRequest) alertsResponse {
	response := requestAlerts(t, alerts, body)

	assert.Equal(t, 200, response.StatusCode)

	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)

	responseBody := alertsResponse{}
	err = json.Unmarshal(data, &responseBody)
	assert.NoError(t, err)

	return responseBody
}

func TestAlertsEndpoint_ServeHTTP(t *testing.T) {
	alerts := make(alert.Alerts)
	alerts["core"] = []alert.Alert{{
		Title:            "Defect",
		Description:      "Incredible Defect",
		IssuedAt:         alert.IssuedAt{Time: time.Now()},
		AffectedVersions: alert.MustParseVersionRange("<2.27.3"),
	}}

	body := alertsRequest{
		InstanceId: "42",
		Version:    alert.MustParseVersion("2.27.1"),
		Os:         "Linux",
		Arch:       "amd64",
		Java:       "1.8.0_121",
	}

	responseBody := callAlertsEndpoint(t, alerts, body)

	assert.Len(t, responseBody.Alerts, 1)
	assert.Equal(t, responseBody.Alerts[0].Title, "Defect")
	assert.Empty(t, responseBody.Plugins)
}

func TestAlertsEndpoint_ServeHTTPNoAlertsAtAll(t *testing.T) {
	alerts := make(alert.Alerts)

	body := alertsRequest{
		InstanceId: "42",
		Version:    alert.MustParseVersion("2.27.1"),
		Os:         "Linux",
		Arch:       "amd64",
		Java:       "1.8.0_121",
	}

	responseBody := callAlertsEndpoint(t, alerts, body)
	assert.Empty(t, responseBody.Alerts)
	assert.Empty(t, responseBody.Plugins)
}

func TestAlertsEndpoint_ServeHTTPCoreRangeDoesNotMatch(t *testing.T) {
	alerts := make(alert.Alerts)
	alerts["core"] = []alert.Alert{{
		Title:            "Defect",
		Description:      "Incredible Defect",
		IssuedAt:         alert.IssuedAt{Time: time.Now()},
		AffectedVersions: alert.MustParseVersionRange("<2.27.3"),
	}}

	body := alertsRequest{
		InstanceId: "42",
		Version:    alert.MustParseVersion("2.28.0"),
		Os:         "Linux",
		Arch:       "amd64",
		Java:       "1.8.0_121",
	}

	responseBody := callAlertsEndpoint(t, alerts, body)

	assert.Empty(t, responseBody.Alerts)
	assert.Empty(t, responseBody.Plugins)
}

func TestAlertsEndpoint_ServeHTTPPlugins(t *testing.T) {
	alerts := make(alert.Alerts)
	alerts["scm-review-plugin"] = []alert.Alert{{
		Title:            "Defect",
		Description:      "Incredible Defect",
		IssuedAt:         alert.IssuedAt{Time: time.Now()},
		AffectedVersions: alert.MustParseVersionRange(">1.0.0 <2.0.0"),
	}}

	body := alertsRequest{
		InstanceId: "42",
		Version:    alert.MustParseVersion("2.28.0"),
		Os:         "Linux",
		Arch:       "amd64",
		Java:       "1.8.0_121",
		Plugins: []plugin{{
			Name:    "scm-review-plugin",
			Version: alert.MustParseVersion("1.2.1"),
		}},
	}

	responseBody := callAlertsEndpoint(t, alerts, body)

	assert.Empty(t, responseBody.Alerts)
	assert.Len(t, responseBody.Plugins, 1)
	assert.Equal(t, "scm-review-plugin", responseBody.Plugins[0].Name)
	assert.Equal(t, "Defect", responseBody.Plugins[0].Alerts[0].Title)
}

func TestAlertsEndpoint_ServeHTTPWithGet(t *testing.T) {
	endpoint := CreateAlertsEndpoint(make(alert.Alerts))

	request, err := http.NewRequest("GET", "/api/v1/alerts", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	endpoint.ServeHTTP(w, request)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Result().StatusCode)
}

func TestAlertsEndpoint_ServeHTTPWithoutBody(t *testing.T) {
	endpoint := CreateAlertsEndpoint(make(alert.Alerts))

	request, err := http.NewRequest("POST", "/api/v1/alerts", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	endpoint.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestAlertsEndpoint_ServeHTTPWithNonValidJSON(t *testing.T) {
	endpoint := CreateAlertsEndpoint(make(alert.Alerts))

	request, err := http.NewRequest("POST", "/api/v1/alerts", bytes.NewReader([]byte("__")))
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	endpoint.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestAlertsEndpoint_ServeHTTPWithoutInstanceId(t *testing.T) {
	body := alertsRequest{
		InstanceId: "",
		Version:    alert.MustParseVersion("2.0.0"),
		Os:         "Linux",
		Arch:       "amd64",
		Java:       "9",
	}

	response := requestAlerts(t, make(alert.Alerts), body)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestAlertsEndpoint_ServeHTTPWithoutVersion(t *testing.T) {
	body := make(map[string]string)
	body["instanceId"] = "42"
	body["os"] = "Windows"
	body["arch"] = "arm"
	body["java"] = "14"

	response := requestAlerts(t, make(alert.Alerts), body)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestAlertsEndpoint_ServeHTTPWithoutPluginName(t *testing.T) {
	body := alertsRequest{
		InstanceId: "42",
		Version:    alert.MustParseVersion("2.0.0"),
		Os:         "Linux",
		Arch:       "amd64",
		Java:       "9",
		Plugins:    []plugin{{Version: alert.MustParseVersion("1.0.0")}},
	}

	response := requestAlerts(t, make(alert.Alerts), body)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestAlertsEndpoint_ServeHTTPWithoutPluginVersion(t *testing.T) {
	var plugins []map[string]string

	plugin := make(map[string]string)
	plugin["name"] = "scm-review-plugin"

	plugins = append(plugins, plugin)

	body := make(map[string]interface{})
	body["instanceId"] = "42"
	body["version"] = "2.0.0"
	body["os"] = "Windows"
	body["arch"] = "arm"
	body["java"] = "14"
	body["plugins"] = plugins

	response := requestAlerts(t, make(alert.Alerts), body)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestAlertsEndpoint_ServeHTTPCollectMetrics(t *testing.T) {
	alertsRequestCounter.Reset()

	body := alertsRequest{
		InstanceId: "42",
		Version:    alert.MustParseVersion("2.28.0"),
		Os:         "Linux",
		Arch:       "amd64",
		Java:       "1.8.0_121",
		Plugins: []plugin{{
			Name:    "scm-review-plugin",
			Version: alert.MustParseVersion("1.2.1"),
		}},
	}

	callAlertsEndpoint(t, make(alert.Alerts), body)

	counter, err := alertsRequestCounter.GetMetricWithLabelValues("42", alert.CORE, "2.28.0", "Linux", "amd64", "1.8.0_121")
	assert.NoError(t, err)

	v := testutil.ToFloat64(counter)
	assert.Equal(t, 1.0, v)

	counter, err = alertsRequestCounter.GetMetricWithLabelValues("42", "scm-review-plugin", "1.2.1", "Linux", "amd64", "1.8.0_121")
	assert.NoError(t, err)

	v = testutil.ToFloat64(counter)
	assert.Equal(t, 1.0, v)
}
