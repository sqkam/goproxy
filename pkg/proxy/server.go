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

const bufSize = 512

func (s *server) copyHeader(r, req *http.Request) {
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, r.Method, s.target, r.Body)
	if err != nil {
		log.Printf("http.NewRequestWithContext error: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.copyHeader(r, req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("http.DefaultClient.Do error: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	buf := make([]byte, bufSize)
	_, err = io.CopyBuffer(w, readerx.NewLoggerReader(resp.Body), buf)
	if err != nil {
		log.Printf("io.CopyBuffer error: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *server) Run(ctx context.Context) {
	log.Printf(fmt.Sprintf("Starting http gateway server on port %d...", s.listen))

	err := http.ListenAndServe(fmt.Sprintf(":%d", s.listen), s)
	if err != nil {
		panic(err)
	}
}

func NewDefaultServer(conf *config.ProxyConfig) Service {
	return &server{
		listen: conf.Listen,
		target: conf.Target,
	}
}
