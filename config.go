package mailgunexporter

import (
	"fmt"
	"os"
	"strings"
)

var (
	defaultPort        = 9999
	defaultHTTPAddress = fmt.Sprintf(":%d", defaultPort)
	defaultRegion      = "US"
)

// Config contains all config for the mailgun-exporter
type Config struct {
	// MailgunPrivateAPIKey is your Mailgun private API key
	// You can find the Private API Key in your Account Menu, under "Settings":
	// (https://app.mailgun.com/app/account/security)
	MailgunPrivateAPIKey string

	HTTPAddress string

	MailgunRegion string
}

// NewConfigFromEnv creates a new config instance from the environment
func NewConfigFromEnv() *Config {
	config := &Config{}

	config.MailgunPrivateAPIKey = getEnv("MAILGUN_PRIVATE_API_KEY", "")

	config.HTTPAddress = getEnv("HTTP_ADDRESS", defaultHTTPAddress)

	config.MailgunRegion = getEnv("MAILGUN_REGION", defaultRegion)

	return config
}

// Validate validates the config
func (config *Config) Validate() error {

	if config.MailgunPrivateAPIKey == "" {
		return fmt.Errorf("expected private api key in config to be set")
	}

	if config.HTTPAddress == "" {
		return fmt.Errorf("HTTPAddress cannot be empty")
	}

	if strings.ToUpper(config.MailgunRegion) != "EU" && strings.ToUpper(config.MailgunRegion) != "US" {
		return fmt.Errorf("invalid mailgun region: %q", config.MailgunRegion)
	}

	return nil
}

// getEnv gets the env variable with the given key if the key exists
// else it falls back to the fallback value
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
