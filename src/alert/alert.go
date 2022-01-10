package alert

import (
	"encoding/json"
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

func (t *IssuedAt) UnmarshalJSON(data []byte) error {
	return t.UnmarshalYAML(func(i interface{}) error {
		return json.Unmarshal(data, i)
	})
}

func (t *IssuedAt) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.Format("2006-01-02"))
}

type VersionRange struct {
	Value string
	Range semver.Range
}

func MustParseVersionRange(value string) VersionRange {
	r := semver.MustParseRange(value)
	return VersionRange{Value: value, Range: r}
}

func (r *VersionRange) Contains(version string) (bool, error) {
	v, err := semver.Parse(version)
	if err != nil {
		return false, errors.Wrapf(err, "Failed to parse semver %s", version)
	}
	return r.Range(v), nil
}

func (r *VersionRange) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Value)
}

func (r *VersionRange) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var buf string
	err := unmarshal(&buf)
	if err != nil {
		return nil
	}

	ra, err := semver.ParseRange(buf)
	if err != nil {
		return errors.Wrap(err, "Failed to parse range")
	}
	r.Value = buf
	r.Range = ra
	return nil
}

func (r *VersionRange) UnmarshalJSON(data []byte) error {
	return r.UnmarshalYAML(func(i interface{}) error {
		return json.Unmarshal(data, i)
	})
}

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
