package http

import (
	"fmt"
	"time"
)

func GetTime(buf []byte) []byte {
	return time.Now().AppendFormat(buf, "Mon, 02 Jan 2006 15:04:05 GMT")
}

func WriteOKResponse(hc *HttpCodec, contentType string, body []byte) {
	hc.AppendToBuffer([]byte("HTTP/1.1 200 OK\r\nServer: gnet\r\nContent-Type: " + contentType + "\r\nDate: "))
	hc.AppendToBuffer(GetTime(hc.GetBuffer()))
	hc.AppendToBuffer([]byte(fmt.Sprintf("\r\nContent-Length: %d\r\n\r\n", len(body))))
	hc.AppendToBuffer(body)
}

func WriteErrorResponse(hc *HttpCodec, statusCode int, statusText string, message string) {
	hc.AppendToBuffer([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\nServer: gnet\r\nContent-Type: text/plain; charset=utf-8\r\nDate: ", statusCode, statusText)))
	hc.AppendToBuffer(GetTime(hc.GetBuffer()))
	hc.AppendToBuffer([]byte(fmt.Sprintf("\r\nContent-Length: %d\r\n\r\n", len(message))))
	hc.AppendToBuffer([]byte(message))
}

func WriteNotFoundResponse(hc *HttpCodec) {
	WriteErrorResponse(hc, 404, "Not Found", "404 page not found")
}

func WriteInternalServerErrorResponse(hc *HttpCodec, message string) {
	WriteErrorResponse(hc, 500, "Internal Server Error", message)
}
