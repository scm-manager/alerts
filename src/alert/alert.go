package alert

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Alert struct {
	Title            string       `json:"title"`
	Description      string       `json:"description"`
	Link             string       `json:"link"`
	IssuedAt         IssuedAt     `yaml:"issuedAt" json:"issuedAt"`
	AffectedVersions VersionRange `yaml:"affectedVersions" json:"affectedVersions"`
}

func ReadFromFile(filePath string) (Alert, error) {
	alert := Alert{}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return alert, errors.Wrapf(err, "Failed to readAlertsFromDirectory alert from file %s", filePath)
	}

	err = yaml.Unmarshal(data, &alert)
	if err != nil {
		return alert, errors.Wrapf(err, "Failed to unmarshal alert from file %s", filePath)
	}
	return alert, nil
}
