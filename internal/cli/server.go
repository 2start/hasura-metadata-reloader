package cli

import (
	"fmt"
	"net/http"
	"os"

	"github.com/2start/hasura-metadata-reloader/internal/config"
	"github.com/2start/hasura-metadata-reloader/internal/handler"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server --endpoint=https://hasura.canida.io --admin-secret=MYSECRET",
		Short: "Start a server with a single endpoint to call the Hasura API to reload metadata.",
		Long: "Start a server to with a single endpoint to call the Hasura API to reload metadata. If there is an inconsistency found it will be logged." +
			" If sentry-dsn is specified, the inconsistency will be reported to Sentry.",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := config.NewConfigurationWithCmdParams(cmd)
			if err != nil {
				os.Exit(1)
			}
			metaHandler := &handler.MetadataHandler{
				Config: config,
			}

			http.HandleFunc("/reload_metadata", metaHandler.HandleReloadMetadataRequest)

			bindAddress := fmt.Sprintf(":%d", config.Port)

			log.Info().Int("port", config.Port).Msg("Server running..")
			if err := http.ListenAndServe(bindAddress, nil); err != nil {
				log.Fatal().Err(err).Msg("Failed to start the server.")
			}
		},
	}
)

func init() {
	serverCmd.Flags().String("hasura-endpoint", "", "Hasura endpoint URL, e.g. https://hasura.canida.io")
	serverCmd.Flags().String("hasura-admin-secret", "", "Hasura admin secret.")
	rootCmd.AddCommand(serverCmd)
}
