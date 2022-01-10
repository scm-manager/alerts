package alert

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadFromDirectory_FindCoreAlerts(t *testing.T) {
	alerts, err := ReadFromDirectory("testdata/001")
	assert.NoError(t, err)

	coreAlerts := alerts[CORE]
	assert.Len(t, coreAlerts, 1)
	assert.Equal(t, "Logback vulnerability", coreAlerts[0].Title)
}

func TestReadFromDirectory_FindPlugins(t *testing.T) {
	alerts, err := ReadFromDirectory("testdata/001")
	assert.NoError(t, err)

	pluginAlerts := alerts["scm-review-plugin"]
	assert.Len(t, pluginAlerts, 1)
	assert.Equal(t, "Special", pluginAlerts[0].Title)
}

func TestReadFromDirectory_FindPluginWithMultipleAlerts(t *testing.T) {
	alerts, err := ReadFromDirectory("testdata/001")
	assert.NoError(t, err)

	pluginAlerts := alerts["scm-editor-plugin"]
	assert.Len(t, pluginAlerts, 2)
	assert.NotEqual(t, pluginAlerts[0].Title, pluginAlerts[1].Title)
}

func TestReadFromDirectory_PluginWithoutAlertsDir(t *testing.T) {
	alerts, err := ReadFromDirectory("testdata/001")
	assert.NoError(t, err)

	pluginAlerts := alerts["scm-pathwp-plugin"]
	assert.Nil(t, pluginAlerts)
}
