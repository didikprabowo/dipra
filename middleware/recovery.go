package middleware

import (
	"fmt"
	"log"
	"runtime"

	"github.com/didikprabowo/dipra"
)

func Recovery() dipra.MiddlewareFunc {
	return func(next dipra.HandlerFunc) dipra.HandlerFunc {
		return func(c *dipra.Context) (err error) {
			defer func() {
				if errS := recover(); errS != nil {
					err, ok := errS.(error)
					if !ok {
						err = fmt.Errorf("%s", errS)
					}

					_, file, line, okC := runtime.Caller(3)
					if okC {
						log.Printf("[Recovery] Panic Line => %d File => %s, error : %s", line, file, err.Error())
					} else {
						log.Printf("[Recovery] Panic recovery : %s", err.Error())
					}

					c.SetError(err)
				}
			}()

			return next(c)
		}
	}
}
