package alert

import (
	"github.com/blang/semver/v4"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"time"
)

type IssuedAt struct {
	time.Time
}

func (t *IssuedAt) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var buf string
	err := unmarshal(&buf)
	if err != nil {
		return nil
	}

	tt, err := time.Parse("2006-01-02", strings.TrimSpace(buf))
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}

type Range struct {
	Value string
	semver.Range
}

func (r *Range) Contains(version string) (bool, error) {
	v, err := semver.Parse(version)
	if err != nil {
		return false, errors.Wrapf(err, "Failed to parse semver %s", version)
	}
	return r.Range(v), nil
}

func (r *Range) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var buf string
	err := unmarshal(&buf)
	if err != nil {
		return nil
	}

	ra, err := semver.ParseRange(buf)
	if err != nil {
		return errors.Wrap(err, "Failed to parse range")
	}
	r.Range = ra
	return nil
}

type Alert struct {
	Title            string
	Description      string
	Link             string
	IssuedAt         IssuedAt `yaml:"issuedAt"`
	AffectedVersions Range    `yaml:"affectedVersions"`
}

func ReadFromFile(filePath string) (Alert, error) {
	alert := Alert{}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return alert, errors.Wrapf(err, "Failed to read alert from file %s", filePath)
	}

	err = yaml.Unmarshal(data, &alert)
	if err != nil {
		return alert, errors.Wrapf(err, "Failed to unmarshal alert from file %s", filePath)
	}
	return alert, nil
}
