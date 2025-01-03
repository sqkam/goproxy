package readerx

import (
	"io"
	"log"
)

type loggerReader struct {
	reader io.Reader
}

func (r *loggerReader) Read(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	if err == nil || err == io.EOF {
		log.Printf("服务器读到数据:\n  %v\n", string(p))
	}
	return n, err
}

func NewLoggerReader(reader io.Reader) io.Reader {
	return &loggerReader{reader: reader}
}
