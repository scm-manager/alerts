package api

import (
	"encoding/json"
	"github.com/scm-manager/alerts/src/alert"
	"log"
	"net/http"
)

type plugin struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type alertsRequest struct {
	Version string   `json:"version"`
	Os      string   `json:"os"`
	Arch    string   `json:"arch"`
	Java    string   `json:"java"`
	Plugins []plugin `json:"plugins"`
}

type pluginAlerts struct {
	Name   string
	Alerts []alert.Alert
}

type alertsResponse struct {
	Alerts  []alert.Alert
	Plugins []pluginAlerts
}

type AlertsEndpoint struct {
	alerts alert.Alerts
}

func (a *AlertsEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := alertsResponse{
		Alerts: a.alerts[alert.CORE],
	}

	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(&response)
	if err != nil {
		http.Error(w, "Failed to marshal alerts to json", 500)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		log.Println("Failed to write response", err)
	}
}

func CreateAlertsEndpoint(alerts alert.Alerts) *AlertsEndpoint {
	return &AlertsEndpoint{alerts}
}
