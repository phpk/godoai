package serv

import (
	"context"
	"godoai/libs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

// CORS 中间件
func corsMiddleware() mux.MiddlewareFunc {
	allowHeaders := "*"
	allowMethods := "GET, POST, PUT, DELETE, OPTIONS"

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			// 允许任何请求头
			w.Header().Set("Access-Control-Allow-Headers", allowHeaders)
			w.Header().Set("Access-Control-Allow-Methods", allowMethods)
			w.Header().Set("Access-Control-Max-Age", "86400")
			// 如果是预检请求（OPTIONS），直接返回 200 OK
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
func LicenseMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			urlPath := r.URL.Path
			if urlPath == "/setLicense" || urlPath == "/getLicense" || strings.HasPrefix(urlPath, "/static") {
				// 直接放行
				next.ServeHTTP(w, r)
				return
			}
			if libs.CheckLinese() {
				next.ServeHTTP(w, r)
			}
		})
	}
}
func (s *Server) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		//log.Println("the request url:", r.RequestURI)
		if *limit == "" {
			// 启动处理关闭信号的goroutine
			go s.checkShutDown()
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
func (s *Server) excuteLimit() {
	if *limit == "" {
		//defer close(closeSignal)
		// 启动处理关闭信号的goroutine
		go s.handleShutdown()
		// 启动空闲检测goroutine，确保在设置 lastUpdateTime 之前执行
		go func() {
			lastUpdateTime = time.Now().Add(-timeOutNum)
			s.checkShutDown()
		}()
	}
}

func (s *Server) shutDown() {
	if err := s.srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}
}

func (s *Server) handleShutdown() {
	sig := <-closeSignal
	log.Printf("Received signal %s, shutting down...\n", sig)
	s.shutDown()
	os.Exit(0)
}

// checkShutDown 用于检查服务器是否处于空闲状态，并根据空闲时间决定是否触发关闭逻辑。
// 该函数不接受参数，也不返回任何值。
func (s *Server) checkShutDown() {
	// 如果 idleTimer 不为空，则停止现有计时器，防止重复调度。
	if idleTimer != nil {
		idleTimer.Stop()
	}

	// 检查当前时间是否已超过最后更新时间加上超时时间，如果是，则触发关闭逻辑。
	if time.Since(lastUpdateTime) > timeOutNum {
		// 重新设置 idleTimer，当经过 timeOutNum 时间后，发送一个空结构体到 closeSignal 通道，以触发相关关闭逻辑。
		log.Println("execute the idleTimer")
		idleTimer = time.AfterFunc(timeOutNum, func() {
			log.Println("idleTimer timeout")
			closeSignal <- struct{}{}
		})
	} else {
		// 如果没有超过超时时间
		lastUpdateTime = time.Now()
		log.Println("execute the lastUpdateTime")
	}
}
func (s *Server) sigChain() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Printf("Received signal %s, shutting down...\n", sig)
		s.shutDown()
		os.Exit(0)
	}()
}
