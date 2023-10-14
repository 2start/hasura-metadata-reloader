package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Configuration struct {
	Port        int
	Endpoint    string
	AdminSecret string
	SentryDsn   string
	SentryEnv   string
}

// LoadFlagsAndBind binds command line flags to Viper environment variables
func (cfg *Configuration) LoadFlagsAndBind(cmd *cobra.Command) error {

	flagsAndEnv := map[string]string{
		"HASURA_ENDPOINT":     "hasura-endpoint",
		"HASURA_ADMIN_SECRET": "hasura-admin-secret",
		"SENTRY_DSN":          "sentry-dsn",
		"SENTRY_ENV":          "sentry-env",
	}

	for envName, flagName := range flagsAndEnv {
		if err := viper.BindPFlag(envName, cmd.Flags().Lookup(flagName)); err != nil {
			return fmt.Errorf("failed to bind flag %s to environment variable %s: %w", flagName, envName, err)
		}
	}

	return nil
}

// LoadConfigurationFromEnv fetches configuration from environment variables, prioritizing flags if set during invocation
func (cfg *Configuration) LoadConfigurationFromEnv() {
	cfg.Port = viper.GetInt("PORT")
	if cfg.Port == 0 {
		cfg.Port = 8080
	}

	cfg.Endpoint = viper.GetString("HASURA_ENDPOINT")
	cfg.AdminSecret = viper.GetString("HASURA_ADMIN_SECRET")
	cfg.SentryDsn = viper.GetString("SENTRY_DSN")
	cfg.SentryEnv = viper.GetString("SENTRY_ENV")
}

// NewConfigurationWithCmdParams creates a new configuration instance and fetches values from both flags and environment variables
func NewConfigurationWithCmdParams(cmd *cobra.Command) (*Configuration, error) {
	cfg := &Configuration{}
	viper.AutomaticEnv()

	if err := cfg.LoadFlagsAndBind(cmd); err != nil {
		return nil, err
	}

	cfg.LoadConfigurationFromEnv()

	return cfg, nil
}

// NewConfigurationFromEnv creates a new configuration instance and fetches values only from environment variables
func NewConfigurationFromEnv() *Configuration {
	viper.AutomaticEnv()

	cfg := &Configuration{}
	cfg.LoadConfigurationFromEnv()

	return cfg
}
