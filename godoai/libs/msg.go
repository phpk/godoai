package libs

import (
	"encoding/json"
	"net/http"
)

// Common response structure
type APIResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
}

func writeJSONResponse(w http.ResponseWriter, res APIResponse, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}

// HTTPError 返回带有JSON错误消息的HTTP错误
func HTTPError(w http.ResponseWriter, status int, message string) {
	writeJSONResponse(w, APIResponse{Message: message, Code: -1}, status)
}
func Error(w http.ResponseWriter, message string) {
	writeJSONResponse(w, APIResponse{Message: message, Code: -1}, 200)
}
func Success(w http.ResponseWriter, data any, message string) {
	writeJSONResponse(w, APIResponse{Message: message, Data: data, Code: 0}, 200)
}
