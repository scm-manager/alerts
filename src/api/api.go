package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/scm-manager/alerts/src/alert"
	"gopkg.in/validator.v2"
	"io/ioutil"
	"log"
	"net/http"
)

type plugin struct {
	Name    string        `json:"name" validate:"nonzero"`
	Version alert.Version `json:"version" validate:"nondefault"`
}

type alertsRequest struct {
	InstanceId string        `json:"instanceId" validate:"nonzero"`
	Version    alert.Version `json:"version" validate:"nondefault"`
	Os         string        `json:"os"`
	Arch       string        `json:"arch"`
	Jre        string        `json:"jre"`
	Plugins    []plugin      `json:"plugins"`
}

type pluginAlerts struct {
	Name   string        `json:"name"`
	Alerts []alert.Alert `json:"alerts"`
}

type alertsResponse struct {
	Alerts  []alert.Alert  `json:"alerts"`
	Plugins []pluginAlerts `json:"plugins"`
}

type AlertsEndpoint struct {
	alerts    alert.Alerts
	validator *validator.Validator
}

func (ae *AlertsEndpoint) collectAlerts(request alertsRequest, name string, version alert.Version) []alert.Alert {
	alertsRequestCounter.WithLabelValues(
		request.InstanceId,
		name,
		version.String(),
		request.Os,
		request.Arch,
		request.Jre,
	).Inc()

	var alerts []alert.Alert

	for _, a := range ae.alerts[name] {
		if a.AffectedVersions.Contains(version) {
			alerts = append(alerts, a)
		}
	}

	return alerts
}

func (ae *AlertsEndpoint) findAlerts(request alertsRequest) alertsResponse {
	coreAlerts := ae.collectAlerts(request, alert.CORE, request.Version)

	var plugins []pluginAlerts
	for _, p := range request.Plugins {
		alerts := ae.collectAlerts(request, p.Name, p.Version)
		plugins = append(plugins, pluginAlerts{Name: p.Name, Alerts: alerts})
	}

	return alertsResponse{
		Alerts:  coreAlerts,
		Plugins: plugins,
	}
}

func (ae *AlertsEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Body == nil {
		http.Error(w, "Missing request body", http.StatusBadRequest)
		return
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Println("Failed to close request body")
		}
	}()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	requestBody := alertsRequest{}
	err = json.Unmarshal(data, &requestBody)
	if err != nil {
		http.Error(w, "Failed to unmarshal request body", http.StatusBadRequest)
		return
	}

	go ae.sendToStatsApi(err, data)

	err = ae.validator.Validate(requestBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to validate request: %v", err), http.StatusBadRequest)
		return
	}

	response := ae.findAlerts(requestBody)

	w.Header().Set("Content-Type", "application/json")

	data, err = json.Marshal(&response)
	if err != nil {
		http.Error(w, "Failed to marshal alerts to json", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		log.Println("Failed to write response", err)
	}
}

func (ae *AlertsEndpoint) sendToStatsApi(err error, data []byte) {
	_, err = http.Post("http://stats-api/api/v1/alerts", "application/json", bytes.NewReader(data))
	if err != nil {
		log.Printf("Could not send request to stats-api: %s", err)
	}
}

func nonDefaultVersion(v interface{}, _ string) error {
	version, ok := v.(alert.Version)
	if !ok {
		return validator.ErrUnsupported
	}

	if version.IsDefault() {
		return errors.New("Version must not be empty and not 0.0.0")
	}
	return nil
}

func CreateAlertsEndpoint(alerts alert.Alerts) *AlertsEndpoint {
	v := validator.NewValidator()
	err := v.SetValidationFunc("nondefault", nonDefaultVersion)
	if err != nil {
		log.Fatal("Failed to create custom validation func", err)
	}
	return &AlertsEndpoint{alerts: alerts, validator: v}
}
