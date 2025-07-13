package config

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetConfig(t *testing.T) {
	cfg, err := GetConfig()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
}

func TestLoadJSONConfig(t *testing.T) {
	jsonConfig := JSONConfig{
		Address:        "localhost:9090",
		ReportInterval: "5s",
		PollInterval:   "1s",
		CryptoKey:      "/path/to/key.pem",
	}

	data, err := json.Marshal(jsonConfig)
	require.NoError(t, err)

	tmpFile, err := os.CreateTemp("", "config_test_*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(data)
	require.NoError(t, err)
	tmpFile.Close()

	loadedConfig, err := loadJSONConfig(tmpFile.Name())
	require.NoError(t, err)

	assert.Equal(t, "localhost:9090", loadedConfig.Address)
	assert.Equal(t, "5s", loadedConfig.ReportInterval)
	assert.Equal(t, "1s", loadedConfig.PollInterval)
	assert.Equal(t, "/path/to/key.pem", loadedConfig.CryptoKey)
}

func TestApplyJSONConfig(t *testing.T) {
	config := Config{
		Address:        "localhost:8080",
		PollInterval:   2,
		ReportInterval: 10,
		CryptoKey:      "",
	}

	jsonConfig := &JSONConfig{
		Address:        "localhost:9090",
		ReportInterval: "5s",
		PollInterval:   "1s",
		CryptoKey:      "/path/to/key.pem",
	}

	applyJSONConfig(&config, jsonConfig)

	assert.Equal(t, "localhost:9090", config.Address)
	assert.Equal(t, 5, config.ReportInterval)
	assert.Equal(t, 1, config.PollInterval)
	assert.Equal(t, "/path/to/key.pem", config.CryptoKey)
}

func TestApplyJSONConfigPartial(t *testing.T) {
	config := Config{
		Address:        "localhost:8080",
		PollInterval:   2,
		ReportInterval: 10,
		CryptoKey:      "",
	}

	jsonConfig := &JSONConfig{
		Address:   "localhost:9090",
		CryptoKey: "/path/to/key.pem",
	}

	applyJSONConfig(&config, jsonConfig)

	assert.Equal(t, "localhost:9090", config.Address)
	assert.Equal(t, 10, config.ReportInterval)
	assert.Equal(t, 2, config.PollInterval)
	assert.Equal(t, "/path/to/key.pem", config.CryptoKey)
}

func TestApplyJSONConfigWithDuration(t *testing.T) {
	config := Config{
		Address:        "localhost:8080",
		PollInterval:   2,
		ReportInterval: 10,
		CryptoKey:      "",
	}

	jsonConfig := &JSONConfig{
		ReportInterval: "1m",
		PollInterval:   "30s",
	}

	applyJSONConfig(&config, jsonConfig)

	assert.Equal(t, 60, config.ReportInterval)
	assert.Equal(t, 30, config.PollInterval)
}
