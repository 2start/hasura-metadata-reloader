package hasura

import (
	"encoding/json"

	"github.com/2start/hasura-metadata-reloader/internal/reporting"
	"github.com/rs/zerolog/log"
)

type ReloadMetadataRequestBody struct {
	Type string                 `json:"type"`
	Args map[string]interface{} `json:"args"`
}

type ReloadMetadataResponseBody struct {
	IsConsistent        bool            `json:"is_consistent"`
	InconsistentObjects json.RawMessage `json:"inconsistent_objects"`
}

type MetadataInconsistentError struct {
	InconsistentObjects json.RawMessage
}

func (e *MetadataInconsistentError) Error() string {
	return "metadata is inconsistent"
}

func ReloadMetadata(endpoint string, adminSecret string) error {
	metadataClient := NewClient(endpoint, adminSecret)
	jsonResp, err := metadataClient.SendRequest("reload_metadata", make(map[string]interface{}))

	if err != nil {
		log.Error().Err(err).Msg("The HTTP request to the Metadata API failed.")
		return err
	}

	var respBody ReloadMetadataResponseBody
	err = json.Unmarshal(jsonResp, &respBody)

	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to parse the response.")
		return err
	}

	if !respBody.IsConsistent {
		errInconsistency := &MetadataInconsistentError{respBody.InconsistentObjects}

		log.Error().
			Err(errInconsistency).
			RawJSON("inconsistent_objects", respBody.InconsistentObjects).
			Msg("Metadata is inconsistent.")

		indendetJsonInconsistentObjects, _ := json.MarshalIndent(respBody.InconsistentObjects, "", "  ")
		inconsistentObjects := string(indendetJsonInconsistentObjects)
		reporting.CaptureErrorWithContext(
			errInconsistency,
			"hasura",
			map[string]interface{}{
				"endpoint":             endpoint,
				"is_consistent":        respBody.IsConsistent,
				"inconsistent_objects": inconsistentObjects,
			},
		)

		return errInconsistency
	}

	log.Info().
		Bool("is_consistent", respBody.IsConsistent).
		Msg("Metadata is consistent.")

	return nil
}
