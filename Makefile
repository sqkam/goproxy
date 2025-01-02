format:
	gofumpt -l -w .

build:
	make format
	go build -ldflags '-w -s' -trimpath github.com/sqkam/goproxy/cmd/goproxy
