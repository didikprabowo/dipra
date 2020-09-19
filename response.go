package dipra

import (
	"bufio"
	"net"
	"net/http"
)

type (
	// ResponseWriter ...
	ResponseWriter struct {
		Response   http.ResponseWriter
		statusCode int
		Header     Header
	}
)

// Reset ResponseWriter
func (rsp *ResponseWriter) Reset(w http.ResponseWriter) {
	rsp.Response = w
}

// Write Response
func (rsp *ResponseWriter) Write(data []byte) {
	rsp.Response.Write(data)
}

// WriteHeader Key value
func (rsp *ResponseWriter) WriteHeader(v map[string]string) {
	for k, v := range v {
		rsp.Response.Header().Add(k, v)
	}
}

// WriteStatus Status Code
func (rsp *ResponseWriter) WriteStatus(status int) {
	rsp.statusCode = status
	rsp.Response.WriteHeader(status)
}

// Hijack implements the http.Hijacker
func (rsp *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return rsp.Response.(http.Hijacker).Hijack()
}

// CloseNotify implements the http.CloseNotify
func (rsp *ResponseWriter) CloseNotify() <-chan bool {
	return rsp.Response.(http.CloseNotifier).CloseNotify()
}

// Flush http.Flush
func (rsp *ResponseWriter) Flush() {
	rsp.Response.(http.Flusher).Flush()
}

// Pusher http.pusher
func (rsp *ResponseWriter) Pusher() (pusher http.Pusher) {
	if pusher, ok := rsp.Response.(http.Pusher); ok {
		return pusher
	}
	return nil
}
