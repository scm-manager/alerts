package alert

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const YAML_001 = "testdata/logback.yaml"

func TestReadFromFile_StringFields(t *testing.T) {
	alert, err := ReadFromFile(YAML_001)
	assert.NoError(t, err)

	assert.Equal(t, "Logback vulnerability", alert.Title)
	assert.Equal(t, "The logback team announced a vuln in logback versions", alert.Description)
	assert.Equal(t, "https://scm-manager.org/blog/posts/2021-12-13-log4shell/", alert.Link)
}

func TestReadFromFile_IssuedAtFieldAsTime(t *testing.T) {
	alert, err := ReadFromFile(YAML_001)
	assert.NoError(t, err)

	assert.Equal(t, 2021, alert.IssuedAt.Year())
	assert.Equal(t, time.Month(12), alert.IssuedAt.Month())
	assert.Equal(t, 13, alert.IssuedAt.Day())
}

func TestReadFromFile_AffectedVersionsAsRange(t *testing.T) {
	alert, err := ReadFromFile(YAML_001)
	assert.NoError(t, err)

	r := alert.AffectedVersions

	assert.True(t, r.Contains(MustParseVersion("2.27.1")))
	assert.False(t, r.Contains(MustParseVersion("2.28.0")))
}

func TestReadFromFile_NotExisting(t *testing.T) {
	_, err := ReadFromFile("testdata/notfound")
	assert.Error(t, err)
}

func TestReadFromFile_Invalid(t *testing.T) {
	_, err := ReadFromFile("testdata/noyaml")
	assert.Error(t, err)
}

func TestAlert_JSONMarshal(t *testing.T) {
	alert, err := ReadFromFile(YAML_001)
	assert.NoError(t, err)

	data, err := json.Marshal(&alert)
	assert.NoError(t, err)

	nodes := make(map[string]string)
	err = json.Unmarshal(data, &nodes)
	assert.NoError(t, err)

	assert.Equal(t, "Logback vulnerability", nodes["title"])
	assert.Equal(t, "The logback team announced a vuln in logback versions", nodes["description"])
	assert.Equal(t, "https://scm-manager.org/blog/posts/2021-12-13-log4shell/", nodes["link"])
	assert.Equal(t, "2021-12-13", nodes["issuedAt"])
	assert.Equal(t, "<2.27.3", nodes["affectedVersions"])
}

func TestAlert_JSONUnmarshal(t *testing.T) {
	alert, err := ReadFromFile(YAML_001)
	assert.NoError(t, err)

	data, err := json.Marshal(&alert)
	assert.NoError(t, err)

	alert = Alert{}
	err = json.Unmarshal(data, &alert)
	assert.NoError(t, err)

	assert.Equal(t, "Logback vulnerability", alert.Title)
	assert.Equal(t, "The logback team announced a vuln in logback versions", alert.Description)
	assert.Equal(t, "https://scm-manager.org/blog/posts/2021-12-13-log4shell/", alert.Link)
	assert.Equal(t, "2021-12-13", alert.IssuedAt.Format("2006-01-02"))
	assert.Equal(t, "<2.27.3", alert.AffectedVersions.Value)
}
