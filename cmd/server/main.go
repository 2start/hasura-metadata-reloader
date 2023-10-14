package main

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/2start/hasura-metadata-reloader/internal/config"
	"github.com/2start/hasura-metadata-reloader/internal/handler"
)

func main() {
	config := config.NewConfigurationFromEnv()

	metaHandler := &handler.MetadataHandler{
		Config: config,
	}

	http.HandleFunc("/reload_metadata", metaHandler.HandleReloadMetadataRequest)

	bindAddress := fmt.Sprintf(":%d", config.Port)

	log.Info().Int("port", config.Port).Msg("Server starting.")
	if err := http.ListenAndServe(bindAddress, nil); err != nil {
		log.Fatal().Err(err).Msg("Failed to start the server.")
	}
}
