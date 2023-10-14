package cli

import (
	"fmt"
	"os"

	"github.com/2start/hasura-metadata-reloader/internal/config"
	"github.com/2start/hasura-metadata-reloader/internal/reporting"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hasura-metadata-reloader",
	Short: "Hasura Metadata Reloader",
	Long:  "Check hasura metadata for inconsistencies and reload if necessary. Report to Sentry if inconsistency was found.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cfg, err := config.NewConfigurationWithCmdParams(cmd)
		if err != nil {
			os.Exit(1)
		}

		println("endpoint: " + cfg.Endpoint)

		reporting.InitSentry(cfg.SentryDsn, cfg.SentryEnv)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		reporting.FlushSentry()
	},
}

func init() {
	rootCmd.PersistentFlags().String("sentry-dsn", "", "Sentry DSN. If not specified, sentry will not be initialized.")
	rootCmd.PersistentFlags().String("sentry-env", "", "Optionally, set the Sentry Environment to e.g. staging, production, local.")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
