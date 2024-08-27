package embeddings

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"godoai/knowledge/dbtype"
)

//const defaultBaseURLOllama = "http://localhost:8210/api"

type ollamaResponse struct {
	Embedding []float32 `json:"embedding"`
}

// NewEmbeddingFuncOllama directly sends a request to get embeddings and returns them.
func NewEmbeddingFuncOllama(config dbtype.EmbedConifig, model string, text string) ([]float32, error) {
	//baseURLOllama := defaultBaseURLOllama
	client := &http.Client{}

	// Prepare the request body.
	reqBody, err := json.Marshal(map[string]string{
		"model":  model,
		"prompt": text,
		"type":   config.ApiType,
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal request body: %w", err)
	}

	// Create the request.
	req, err := http.NewRequest("POST", config.ApiUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("couldn't create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request.
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("couldn't send request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status.
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error response from the embedding API: " + resp.Status)
	}

	// Read and decode the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read response body: %w", err)
	}
	var embeddingResponse ollamaResponse
	err = json.Unmarshal(body, &embeddingResponse)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal response body: %w", err)
	}

	// Return the embedding directly.
	return embeddingResponse.Embedding, nil
}
