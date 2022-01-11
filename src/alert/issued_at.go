package alert

import (
	"encoding/json"
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
