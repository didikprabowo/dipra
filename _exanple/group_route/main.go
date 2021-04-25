package main

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/didikprabowo/dipra"
)

func main() {

	r := dipra.Default()
	r.GET("/", func(c *dipra.Context) error {
		return c.String(200, "welcome")
	})
	// Normal Group
	v1 := r.Group("/v1")

	{
		// {basepath}/v1
		v1.GET("/", func(c *dipra.Context) error {
			return c.String(200, fmt.Sprintf("Welcome to api version v1 %+f\n", rand.Float64()))
		})

	}

	// Group use middleware
	v2 := r.Group("/v2", func(hf dipra.HandlerFunc) dipra.HandlerFunc {
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
			return c.String(200, fmt.Sprintf("Welcome to api version v2 %+f\n", rand.Float64()))
		})

	}

	err := r.Run(":9020")
	if err != nil {
		panic(err)
	}
}
