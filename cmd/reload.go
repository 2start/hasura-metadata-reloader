package cmd

import (
	"os"

	"github.com/2start/hasura-metadata-reloader/internal/hasura"
	"github.com/spf13/cobra"
)

var (
	endpoint    string
	adminSecret string
	reloadCmd   = &cobra.Command{
		Use:   "reload --endpoint=https://hasura.canida.io --admin-secret=MYSECRET",
		Short: "Call Hasura API to reload metadata.",
		Long: "Call Hasura API to reload metadata. If there is an inconsistency found it will be logged." +
		" If sentry-dsn is specified, the inconsistency will be reported to Sentry.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := hasura.ReloadMetadata(endpoint, adminSecret); err != nil {
				os.Exit(1)
			}
		},
	}
)

func init() {
	reloadCmd.Flags().StringVar(&endpoint, "endpoint", "",
		"deployment environment for API server, e.g. local, prod (it should match your configuration filename)")
	_ = reloadCmd.MarkFlagRequired("endpoint")
	reloadCmd.Flags().StringVar(&adminSecret, "admin-secret", "",
		"deployment environment for API server, e.g. local, prod (it should match your configuration filename)")
	_ = reloadCmd.MarkFlagRequired("admin-secret")
	rootCmd.AddCommand(reloadCmd)
}
