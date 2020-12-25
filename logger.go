package dipra

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	// Reset color
	Reset = "\033[0m"
	// Red color
	Red = "\033[31m"
	// Green color
	Green = "\033[32m"
	// Yellow color
	Yellow = "\033[33m"
	// Blue color
	Blue = "\033[34m"
	// Purple color
	Purple = "\033[35m"
	// Cyan color
	Cyan = "\033[36m"
	// Gray color
	Gray = "\033[37m"
	// White color
	White = "\033[97m"
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
func Logger() MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) (err error) {

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
			l.StatusCode = c.Writen.statusCode
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
		White, time.Now().Format("02/01/2006 15:04:05"), Reset,
		l.GetColor(), l.Request.Method, Reset,
		l.GetColorStatusCode(l.StatusCode), l.StatusCode, Reset,
		Cyan, l.Request.Host, l.Request.URL, Reset,
		Green, l.Latency, Reset,
	)
	if l.Err != "" {
		out += fmt.Sprintf("| %s%v%s", Red, l.Err, Reset)
	}

	buf.WriteString(out + "\n")
	l.Type.Write(buf.Bytes())
	l.Pool.Put(buf)
}

// GetColor is used get color method
func (l *LoggerConfig) GetColor() string {
	switch l.Request.Method {
	case http.MethodGet:
		return Green
	case http.MethodPost:
		return Blue
	case http.MethodDelete:
		return Red
	case http.MethodPut:
		return Cyan
	case http.MethodPatch:
		return Purple
	case http.MethodOptions:
		return Yellow
	default:
		return White
	}
}

// GetColorStatusCode is used code color status code
func (l *LoggerConfig) GetColorStatusCode(code int) string {
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return Green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return White
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return Yellow
	default:
		return Red
	}
}
