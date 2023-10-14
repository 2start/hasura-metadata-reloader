package handler

import (
	"encoding/json"
	"net/http"

	"github.com/2start/hasura-metadata-reloader/internal/config"
	"github.com/2start/hasura-metadata-reloader/internal/hasura"
)

type MetadataHandler struct {
	Config *config.Configuration
}

func (c *MetadataHandler) HandleReloadMetadataRequest(w http.ResponseWriter, r *http.Request) {
	resp, err := hasura.ReloadMetadata(c.Config.Endpoint, c.Config.AdminSecret)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
