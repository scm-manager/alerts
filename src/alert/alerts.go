package alert

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const CORE = "core"

type Alerts map[string][]Alert

func ReadFromDirectory(directoryPath string) (Alerts, error) {
	alerts := make(Alerts)

	coreAlerts, err := readAlertsFromDirectory(directoryPath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to readAlertsFromDirectory core alerts")
	}
	alerts[CORE] = coreAlerts

	pluginDirectory := path.Join(directoryPath, "plugins")
	plugins, err := ioutil.ReadDir(pluginDirectory)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to list plugin directory %s", pluginDirectory)
	}

	for _, plugin := range plugins {
		directory := path.Join(pluginDirectory, plugin.Name())

		s, err := os.Stat(directory)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to read stats from file %s", directory)
		}

		if s.IsDir() {
			pluginAlerts, err := readAlertsFromDirectory(directory)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to readAlertsFromDirectory %s alerts", plugin.Name())
			}
			if pluginAlerts != nil {
				alerts[plugin.Name()] = pluginAlerts
			}
		}
	}

	return alerts, nil
}

func readAlertsFromDirectory(directory string) ([]Alert, error) {
	alertDirectory := path.Join(directory, "alerts")

	if _, err := os.Stat(alertDirectory); os.IsNotExist(err) {
		return nil, nil
	}

	files, err := ioutil.ReadDir(alertDirectory)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to readAlertsFromDirectory directory %s", directory)
	}

	var alerts []Alert
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".yaml") {
			alertPath := path.Join(alertDirectory, f.Name())
			alert, err := ReadFromFile(alertPath)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to readAlertsFromDirectory alert %s", alertPath)
			}

			alerts = append(alerts, alert)
		}
	}
	return alerts, nil
}
