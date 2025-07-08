package http

import (
	"bytes"
	"errors"
	"io"
	"strconv"

	"github.com/evanphx/wildcat"
	"github.com/panjf2000/gnet/v2"
)

var contentLengthKey = []byte("Content-Length")
var CRLF = []byte("\r\n\r\n")

type HttpCodec struct {
	parser        *wildcat.HTTPParser
	contentLength int
	buf           []byte
	body          []byte
}

func NewHttpCodec() *HttpCodec {
	return &HttpCodec{
		parser:        wildcat.NewHTTPParser(),
		contentLength: -1,
	}
}

func (hc *HttpCodec) Parse(data []byte, c gnet.Conn) (int, error) {
	bodyOffset, err := hc.parser.Parse(data)
	if err != nil {
		return 0, err
	}

	if hc.parser.Post() {
		body, err := io.ReadAll(hc.parser.BodyReader(data[bodyOffset:], c))
		if err != nil {
			return 0, err
		}
		hc.body = body
	}

	contentLength := hc.getContentLength()
	if contentLength > -1 {
		return bodyOffset + contentLength, nil
	}

	if idx := bytes.Index(data, CRLF); idx != -1 {
		return idx + 4, nil
	}

	return 0, errors.New("invalid http request")
}

func (hc *HttpCodec) getContentLength() int {
	if hc.contentLength != -1 {
		return hc.contentLength
	}

	val := hc.parser.FindHeader(contentLengthKey)
	if val != nil {
		i, err := strconv.ParseInt(string(val), 10, 0)
		if err == nil {
			hc.contentLength = int(i)
		}
	}

	return hc.contentLength
}

func (hc *HttpCodec) ResetParser() {
	hc.contentLength = -1
}

func (hc *HttpCodec) Reset() {
	hc.ResetParser()
	hc.buf = hc.buf[:0]
}

func (hc *HttpCodec) GetParser() *wildcat.HTTPParser {
	return hc.parser
}

func (hc *HttpCodec) GetBuffer() []byte {
	return hc.buf
}

func (hc *HttpCodec) AppendToBuffer(data []byte) {
	hc.buf = append(hc.buf, data...)
}

func (hc *HttpCodec) GetBody() []byte {
	return hc.body
}
