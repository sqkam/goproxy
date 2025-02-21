package proxy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/sqkam/goproxy/config"
	"github.com/sqkam/goproxy/pkg/readerx"
)

type server struct {
	listen int64
	target string
}

const (
	defaultTimeout = 60 * time.Second
	bufSize        = 32 * 1024 // 增加buffer大小到32KB以提升性能
)

func (s *server) copyHeader(dst, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 使用父上下文创建超时上下文
	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeout)
	defer cancel()

	// 构建代理请求
	targetURL := s.target + r.URL.String()
	req, err := http.NewRequestWithContext(ctx, r.Method, targetURL, r.Body)
	if err != nil {
		log.Printf("http.NewRequestWithContext error: %v\n", err)
		http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
		return
	}

	// 复制请求头
	s.copyHeader(req.Header, r.Header)
	//
	//// 设置代理相关的头部
	//req.Host = r.Host
	//if clientIP, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
	//	req.Header.Set("X-Forwarded-For", clientIP)
	//}
	//req.Header.Set("X-Forwarded-Host", r.Host)
	//req.Header.Set("X-Forwarded-Proto", "http")

	// 发送代理请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("http.DefaultClient.Do error: %v\n", err)
		http.Error(w, "Failed to proxy request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 复制响应头
	s.copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)

	// 复制响应体
	buf := make([]byte, bufSize)
	_, err = io.CopyBuffer(w, readerx.NewLoggerReader(resp.Body), buf)
	if err != nil {
		log.Printf("io.CopyBuffer error: %v\n", err)
		// 此时已经发送了响应头，无法再修改状态码
		return
	}
}

func (s *server) Run(ctx context.Context) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.listen),
		Handler: s,
	}

	// 在goroutine中启动服务器
	go func() {
		log.Printf("Starting http gateway server on port %d...", s.listen)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// 监听context取消
	<-ctx.Done()

	// 优雅关闭
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server Shutdown error: %v", err)
	}
}

func NewDefaultServer(conf *config.ProxyConfig) Service {
	return &server{
		listen: conf.Listen,
		target: conf.Target,
	}
}
