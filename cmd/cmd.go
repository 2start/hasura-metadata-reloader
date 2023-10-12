package cmd

import (
	"fmt"
	"os"

	"github.com/2start/hasura-metadata-reloader/internal/reporting"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hasura-metadata-reloader",
	Short: "Hasura Metadata Reloader",
	Long:  "Check hasura metadata for inconsistencies and reload if necessary. Report to Sentry if inconsistency was found.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		sentryDsn, _ := cmd.Flags().GetString("sentry-dsn")
		sentryEnv, _ := cmd.Flags().GetString("sentry-env")

		// Initialize logger with the log level
		reporting.InitSentry(sentryDsn, sentryEnv)
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
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
