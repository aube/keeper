package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// EnvConfig holds all configuration parameters for the application.
// Fields are tagged to support multiple configuration sources.
type EnvConfig struct {
	Username              string `mapstructure:"username"`                                              // Server address to listen on
	Password              string `mapstructure:"password"`                                              // Server address to listen on
	ServerAddress         string `mapstructure:"server_address" env:"SERVER_ADDRESS"`                   // Server address to listen on
	FileStoragePath       string `mapstructure:"file_storage_path" env:"FILE_STORAGE_PATH"`             // Path to file storage
	DefaultRequestTimeout int    `mapstructure:"default_request_timeout" env:"DEFAULT_REQUEST_TIMEOUT"` // Default request timeout in seconds
	PublicCertFile        string `mapstructure:"public_cert_file" env:"PUBLIC_CERT_FILE"`
	PrivateCertFile       string `mapstructure:"private_cert_file" env:"PRIVATE_CERT_FILE"`
}

// config() initializes and returns the application configuration.
// It loads configuration in the following order of precedence:
// 1. Command-line flags (highest priority)
// 2. Environment variables
// 3. Configuration file
// 4. Default values (lowest priority)
//
// The configuration is loaded only once, subsequent calls return the cached configuration.
func NewConfig() EnvConfig {
	var config EnvConfig

	// Set default values
	viper.SetDefault("server_address", "localhost:8080")
	viper.SetDefault("token_secret_string", "~_^")
	viper.SetDefault("log_level", "info")
	viper.SetDefault("file_storage_path", "./_storage")
	viper.SetDefault("default_request_timeout", 15)

	viper.AutomaticEnv()

	// Define and parse command-line flags
	pflag.StringP("username", "u", "", "Username")
	pflag.StringP("password", "p", "", "Password")
	pflag.StringP("server_address", "a", "", "Server address to listen on")
	pflag.StringP("file_storage_path", "f", "", "Path to file storage")
	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine) // Flags override everything

	// Unmarshal into struct
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %w", err))
	}

	return config
}
