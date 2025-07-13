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
		Address:       "localhost:9090",
		Restore:       boolPtr(false),
		StoreInterval: "5m",
		StoreFile:     "/tmp/test.json",
		DatabaseDSN:   "postgres://test",
		CryptoKey:     "/path/to/key.pem",
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
	assert.Equal(t, false, *loadedConfig.Restore)
	assert.Equal(t, "5m", loadedConfig.StoreInterval)
	assert.Equal(t, "/tmp/test.json", loadedConfig.StoreFile)
	assert.Equal(t, "postgres://test", loadedConfig.DatabaseDSN)
	assert.Equal(t, "/path/to/key.pem", loadedConfig.CryptoKey)
}

func TestApplyJSONConfig(t *testing.T) {
	config := Config{
		Address:         "localhost:8080",
		StoreInterval:   300,
		FileStoragePath: "./storage.json",
		Restore:         true,
		DatabaseDSN:     "",
		CryptoKey:       "",
	}

	jsonConfig := &JSONConfig{
		Address:       "localhost:9090",
		Restore:       boolPtr(false),
		StoreInterval: "1m",
		StoreFile:     "/tmp/test.json",
		DatabaseDSN:   "postgres://test",
		CryptoKey:     "/path/to/key.pem",
	}

	applyJSONConfig(&config, jsonConfig)

	assert.Equal(t, "localhost:9090", config.Address)
	assert.Equal(t, false, config.Restore)
	assert.Equal(t, 60, config.StoreInterval)
	assert.Equal(t, "/tmp/test.json", config.FileStoragePath)
	assert.Equal(t, "postgres://test", config.DatabaseDSN)
	assert.Equal(t, "/path/to/key.pem", config.CryptoKey)
}

func TestApplyJSONConfigPartial(t *testing.T) {
	config := Config{
		Address:         "localhost:8080",
		StoreInterval:   300,
		FileStoragePath: "./storage.json",
		Restore:         true,
		DatabaseDSN:     "",
		CryptoKey:       "",
	}

	jsonConfig := &JSONConfig{
		Address:     "localhost:9090",
		StoreFile:   "/tmp/test.json",
		DatabaseDSN: "postgres://test",
	}

	applyJSONConfig(&config, jsonConfig)

	assert.Equal(t, "localhost:9090", config.Address)
	assert.Equal(t, true, config.Restore)
	assert.Equal(t, 300, config.StoreInterval)
	assert.Equal(t, "/tmp/test.json", config.FileStoragePath)
	assert.Equal(t, "postgres://test", config.DatabaseDSN)
	assert.Equal(t, "", config.CryptoKey)
}

func boolPtr(b bool) *bool {
	return &b
}
