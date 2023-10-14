package cli

import (
	"os"

	"github.com/2start/hasura-metadata-reloader/internal/config"
	"github.com/2start/hasura-metadata-reloader/internal/hasura"
	"github.com/spf13/cobra"
)

var (
	reloadCmd = &cobra.Command{
		Use:   "reload --endpoint=https://hasura.canida.io --admin-secret=MYSECRET",
		Short: "Call Hasura API to reload metadata.",
		Long: "Call Hasura API to reload metadata. If there is an inconsistency found it will be logged." +
			" If sentry-dsn is specified, the inconsistency will be reported to Sentry.",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.NewConfigurationWithCmdParams(cmd)
			if err != nil {
				os.Exit(1)
			}

			if err := hasura.ReloadMetadata(cfg.Endpoint, cfg.AdminSecret); err != nil {
				os.Exit(1)
			}
		},
	}
)

func init() {
	reloadCmd.Flags().String("endpoint", "", "Hasura endpoint URL, e.g. https://hasura.canida.io")
	reloadCmd.Flags().String("admin-secret", "", "Hasura admin secret.")
	rootCmd.AddCommand(reloadCmd)
}
