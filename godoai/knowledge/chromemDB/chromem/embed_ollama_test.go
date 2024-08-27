package chromem

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"godoai/knowledge/dbtype"
)

func TestNewEmbeddingFuncOllama(t *testing.T) {
	model := "model-small"
	baseURLSuffix := "/api"
	prompt := "hello world"

	wantBody, err := json.Marshal(map[string]string{
		"model":  model,
		"prompt": prompt,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	wantRes := []float32{-0.40824828, 0.40824828, 0.81649655} // normalized version of `{-0.1, 0.1, 0.2}`

	// Mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check URL
		if !strings.HasSuffix(r.URL.Path, baseURLSuffix+"/embeddings") {
			t.Fatal("expected URL", baseURLSuffix+"/embeddings", "got", r.URL.Path)
		}
		// Check method
		if r.Method != "POST" {
			t.Fatal("expected method POST, got", r.Method)
		}
		// Check headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatal("expected Content-Type header", "application/json", "got", r.Header.Get("Content-Type"))
		}
		// Check body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
		if !bytes.Equal(body, wantBody) {
			t.Fatal("expected body", wantBody, "got", body)
		}

		// Write response
		resp := ollamaResponse{
			Embedding: wantRes,
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	// Get port from URL

	config := dbtype.EmbedConifig{
		ApiUrl:  "http://localhost:56711",
		ApiType: "llmgpu",
	}
	f := NewEmbeddingFuncOllama(model, config)
	res, err := f(context.Background(), prompt)
	if err != nil {
		t.Fatal("expected nil, got", err)
	}
	if slices.Compare(wantRes, res) != 0 {
		t.Fatal("expected res", wantRes, "got", res)
	}
}
