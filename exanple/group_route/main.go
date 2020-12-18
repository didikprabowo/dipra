package main

import (
	"net/http"

	"github.com/didikprabowo/dipra"
)

func main() {

	r := dipra.Default()

	// Normal Group
	v1 := r.Group("/v2")

	{
		// {basepath}/v1
		v1.GET("/", func(c *dipra.Context) error {
			return c.String(200, "Welcome to api version v2")
		})

	}

	// Group use middleware
	v2 := r.Group("/v1", func(hf dipra.HandlerFunc) dipra.HandlerFunc {
		return func(c *dipra.Context) (err error) {
			if c.Query("name") == "jhon" {
				return c.String(http.StatusUnprocessableEntity, "Can't continue")
			}
			hf(c)
			return
		}
	})

	{
		// {basepath}/v2
		v2.GET("/", func(c *dipra.Context) error {
			return c.String(200, "Welcome to api version v1")
		})

	}

	err := r.Run(":9020")
	if err != nil {
		panic(err)
	}
}
