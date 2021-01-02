package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/didikprabowo/dipra"
)

type (
	// LoggerConfig is used setting log
	LoggerConfig struct {
		// Set type output , any stdrout and stdrerr
		Type io.Writer
		// Http request
		Request *http.Request
		// Http Response
		Response http.ResponseWriter
		// instance enginer
		StatusCode int
		// Sync
		Pool *sync.Pool
		// Latency
		Latency time.Duration

		// Error
		Err string
	}
)

// DefaultLoggerConfig  is used setting log
func DefaultLoggerConfig() LoggerConfig {
	l := LoggerConfig{
		Type: os.Stdout,
	}
	return l
}

// Logger is used instance log
func Logger() dipra.MiddlewareFunc {
	return func(h dipra.HandlerFunc) dipra.HandlerFunc {
		return func(c *dipra.Context) (err error) {

			l := DefaultLoggerConfig()
			l.Request = c.GetRequest()
			l.Response = c.GetResponse()
			start := time.Now()
			if err := h(c); err != nil {
				c.SetError(err)
				l.Err = err.Error()
			}

			end := time.Now()
			l.Latency = end.Sub(start)
			l.StatusCode = c.Writen.StatusCode
			l.Pool = &sync.Pool{
				New: func() interface{} {
					return bytes.NewBuffer(make([]byte, 256))
				},
			}

			l.BuildLogger()

			return err
		}
	}
}

// BuildLogger is used running log
func (l LoggerConfig) BuildLogger() {
	buf := l.Pool.Get().(*bytes.Buffer)
	buf.Reset()

	out := fmt.Sprintf(
		" [DIPRA] %s%s%s => %s%s%s | %s%d%s | %s\"%s %s\" |%s %s%s %s",
		dipra.White,
		time.Now().Format("02/01/2006 15:04:05"),
		dipra.Reset,
		l.GetColor(),
		l.Request.Method,
		dipra.Reset,
		l.GetColorStatusCode(l.StatusCode),
		l.StatusCode,
		dipra.Reset,
		dipra.Cyan,
		l.Request.Host,
		l.Request.URL,
		dipra.Reset,
		dipra.Green,
		l.Latency,
		dipra.Reset,
	)
	if l.Err != "" {
		out += fmt.Sprintf("| %s%v%s", dipra.Red, l.Err, dipra.Reset)
	}

	buf.WriteString(out + "\n")
	l.Type.Write(buf.Bytes())
	l.Pool.Put(buf)
}

// GetColor is used get color method
func (l *LoggerConfig) GetColor() string {
	switch l.Request.Method {
	case http.MethodGet:
		return dipra.Green
	case http.MethodPost:
		return dipra.Blue
	case http.MethodDelete:
		return dipra.Red
	case http.MethodPut:
		return dipra.Cyan
	case http.MethodPatch:
		return dipra.Purple
	case http.MethodOptions:
		return dipra.Yellow
	default:
		return dipra.White
	}
}

// GetColorStatusCode is used code color status code
func (l *LoggerConfig) GetColorStatusCode(code int) string {
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return dipra.Green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return dipra.White
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return dipra.Yellow
	default:
		return dipra.Red
	}
}
