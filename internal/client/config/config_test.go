package config

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected EnvConfig
	}{
		{
			name: "default values",
			expected: EnvConfig{
				ServerAddress:         "http://localhost:8080/api/v1",
				StoragePath:           "./_storage",
				DefaultRequestTimeout: 15,
				LogLevel:              "info",
			},
		},
		{
			name: "environment variables override",
			envVars: map[string]string{
				"SERVER_ADDRESS":          "http://test:8080",
				"STORAGE_PATH":            "/tmp/storage",
				"DEFAULT_REQUEST_TIMEOUT": "30",
				"LOG_LEVEL":               "debug",
			},
			expected: EnvConfig{
				ServerAddress:         "http://test:8080",
				StoragePath:           "/tmp/storage",
				DefaultRequestTimeout: 30,
				LogLevel:              "debug",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables if any
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			// Reset viper to avoid test pollution
			viper.Reset()
			pflag.CommandLine = pflag.NewFlagSet("test", pflag.ContinueOnError)

			config := NewConfig()

			assert.Equal(t, tt.expected.ServerAddress, config.ServerAddress)
			assert.Equal(t, tt.expected.StoragePath, config.StoragePath)
			assert.Equal(t, tt.expected.DefaultRequestTimeout, config.DefaultRequestTimeout)
			assert.Equal(t, tt.expected.LogLevel, config.LogLevel)
		})
	}
}
