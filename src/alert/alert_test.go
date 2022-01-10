package alert

import (
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

	c, err := r.Contains("2.27.1")
	assert.NoError(t, err)
	assert.True(t, c)

	c, err = r.Contains("2.28.0")
	assert.NoError(t, err)
	assert.False(t, c)
}

func TestReadFromFile_NotExisting(t *testing.T) {
	_, err := ReadFromFile("testdata/notfound")
	assert.Error(t, err)
}

func TestReadFromFile_Invalid(t *testing.T) {
	_, err := ReadFromFile("testdata/noyaml")
	assert.Error(t, err)
}

/**

title: Logback vulnerability
description: The logback team announced a vuln in logback versions
link: https://scm-manager.org/blog/posts/2021-12-13-log4shell/
issuedAt: 2021-12-13
affectedVersions: <2.27.3


*/
