package main

import (
	"net/http"

	"github.com/didikprabowo/dipra"
)

func main() {

	route := dipra.Default()
	s := route.Group("/v1/:base_path")
	s.Any("/*", func(c *dipra.Context) (err error) {
		return c.JSON(http.StatusOK, dipra.M{
			"base_path": c.Param("base_path"),
		})
	})

	// route

	route.Run(":9020")
}
