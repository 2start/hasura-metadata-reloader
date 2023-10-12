package hasura

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

type RequestBody struct {
	Type string                 `json:"type"`
	Args map[string]interface{} `json:"args"`
}

type MetadataClient struct {
	Endpoint    string
	AdminSecret string
	HTTPClient  *http.Client
}

func NewClient(endpoint string, adminSecret string) *MetadataClient {
	return &MetadataClient{
		Endpoint:    endpoint,
		AdminSecret: adminSecret,
		HTTPClient:  &http.Client{},
	}
}

func (c *MetadataClient) SendRequest(reqType string, args map[string]interface{}) (json.RawMessage, error) {
	url := fmt.Sprintf("%s/v1/metadata", c.Endpoint)

	reqBody := &RequestBody{
		Type: reqType,
		Args: args,
	}

	jsonReqBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonReqBody))
	req.Header.Set("X-Hasura-Admin-Secret", c.AdminSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("The HTTP request failed.")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("the HTTP request failed with status code %d", resp.StatusCode)
		return nil, err
	}

	data, _ := io.ReadAll(resp.Body)
	jsonResp := json.RawMessage(data)

	return jsonResp, nil
}
