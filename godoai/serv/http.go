package serv

import (
	"context"
	"flag"
	"godoai/deps"
	"godoai/progress"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/mux"
)

var (
	idleTimer      *time.Timer        // 用于触发超时的定时器。
	timeOutNum     = 30 * time.Minute // 超时时间设置为30分钟
	systemPort     string
	limit          = flag.String("limit", "", "`limit` of concurrent requests")
	closeSignal    = make(chan struct{}, 1) // 用于发送关闭信号的通道，容量为1以避免阻塞。
	lastUpdateTime time.Time
)

// Server 实现了服务器启动和关闭的方法
type Server struct {
	srv    *http.Server
	router *mux.Router
}

func GetServer(port string) *Server {
	systemPort = port
	return newServer()
}

// serv 包内修改 RestrictedRouter
type RestrictedRouter struct {
	router *mux.Router
}

// Handle 注册处理函数
func (rr *RestrictedRouter) Handle(pattern string, method string, handler http.HandlerFunc) {
	rr.router.HandleFunc(pattern, handler).Methods(method)
}

// RegisterRoutes 接收一个函数作为参数，该函数可以注册额外的路由
func (s *Server) RegisterRoutes(registerFunc func(rr *RestrictedRouter)) {
	registerFunc(&RestrictedRouter{s.router})
}

func newServer() *Server {
	router := mux.NewRouter()
	srv := &http.Server{
		Addr:    systemPort,
		Handler: router,
	}
	return &Server{
		router: router,
		srv:    srv,
	}
}

func (s *Server) setupRoutes() {
	s.router.Use(corsMiddleware())
	//s.router.Use(LicenseMiddleware())
	s.router.HandleFunc("/ping", PingHandle).Methods("GET")
	s.router.HandleFunc("/restart", s.RestartHandle).Methods("GET")
	s.router.HandleFunc("/getLicense", GetLicenseHandle).Methods("GET")
	s.router.HandleFunc("/setLicense", SetLicenseHandler).Methods("POST")
	// 使用 http.FileServer 提供静态文件
	distFS, _ := fs.Sub(deps.Frontendassets, "dist")
	fileServer := http.FileServer(http.FS(distFS))
	s.router.PathPrefix("/").Handler(fileServer)

}

func (s *Server) Run() {
	defer close(closeSignal)
	s.setupRoutes()
	s.router.Use(s.middleware)

	s.excuteLimit()

	go func() {
		if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("server failed to start: %v\n", err)
		}
	}()

	s.sigChain()

	log.Printf("Listening on port: %v", s.srv.Addr)
	// 阻塞主goroutine，保持程序运行
	select {}

}

func (s *Server) Start() {
	s.setupRoutes()
	go func() {
		if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("Server failed to start: %v\n", err)
		}
	}()

	s.sigChain()

	log.Printf("Listening on port: %v", s.srv.Addr)
	// 阻塞主goroutine，保持程序运行
	select {}

}
func (s *Server) Stop(ctx context.Context) error {
	defer close(closeSignal)
	return s.srv.Shutdown(ctx)
}
func (s *Server) Restart() error {
	// 获取当前可执行文件的路径
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	// 创建一个新的进程来重启服务
	cmd := exec.Command(exePath)
	cmd = progress.SetHideConsoleCursor(cmd)

	// 启动新进程
	if err := cmd.Start(); err != nil {
		return err
	}
	os.Exit(0)

	return nil
}
