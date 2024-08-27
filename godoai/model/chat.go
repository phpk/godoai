package model

import (
	"encoding/json"
	"godoai/libs"
	"net/http"
)

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	url := libs.GetOllamaUrl() + "/v1/chat/completions"
	var request interface{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		libs.Error(w, err.Error())
		return
	}
	ForwardHandler(w, r, request, url, "POST")
}
func EmbeddingHandler(w http.ResponseWriter, r *http.Request) {
	url := libs.GetOllamaUrl() + "/api/embeddings"
	var request interface{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		libs.Error(w, err.Error())
		return
	}
	ForwardHandler(w, r, request, url, "POST")
}
