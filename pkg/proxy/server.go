package proxy

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/sqkam/goproxy/config"

	"golang.org/x/sync/errgroup"
)

type server struct {
	listen int64
	target string
}

var onceClose sync.Once

func (s *server) Run(ctx context.Context) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.listen))
	if err != nil {
		panic("connection error:" + err.Error())
	}
	defer func() {
		onceClose.Do(func() {
			listener.Close()
		})
	}()

	fmt.Printf(fmt.Sprintf("Starting tcp gateway server on port %d to %s ...\n", s.listen, s.target))
	go func() {
		// context超时 取消代理
		<-ctx.Done()
		onceClose.Do(func() {
			listener.Close()
		})
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept Error:", err)
			continue
		}

		go s.handleTCP(conn)
	}
}

const bufSize = 512

func (s *server) handleTCP(src net.Conn) {
	if src, ok := src.(*net.TCPConn); ok {
		_ = src.SetKeepAlive(true)
		_ = src.SetKeepAlivePeriod(5 * time.Second)
	}

	dst, err := net.DialTimeout("tcp", s.target, time.Second*10)
	if err != nil {
		return
	}

	defer dst.Close()

	var eg errgroup.Group
	eg.Go(func() error {
		cpyBuf := make([]byte, bufSize)
		defer dst.Close()
		_, err := io.CopyBuffer(src, dst, cpyBuf)
		return err
	})
	eg.Go(func() error {
		cpyBuf := make([]byte, bufSize)
		defer src.Close()
		_, err := io.CopyBuffer(dst, src, cpyBuf)
		return err
	})
	err = eg.Wait()
	if err != nil {
		fmt.Printf("handleTCP error %v\n", err.Error())
	}
}

func NewDefaultServer(conf *config.ProxyConfig) Service {
	return &server{
		listen: conf.Listen,
		target: conf.Target,
	}
}
