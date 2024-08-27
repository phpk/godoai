package serv

import (
	"fmt"
	"godoai/libs"
	"net/http"
)

// pingHandler 检查服务是否运行正常
func PingHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK) // 返回200状态码
	libs.Success(w, "", "Pong! Service is running.")
}
func HomeHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK) // 返回200状态码
	fmt.Fprintln(w, "Service is running.")
}
func (s *Server) RestartHandle(w http.ResponseWriter, r *http.Request) {
	// 尝试重启服务
	if err := s.Restart(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to restart service: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Service restarted")
}
