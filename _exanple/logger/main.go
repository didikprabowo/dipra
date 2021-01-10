package main

import (
	"net/http"

	"github.com/didikprabowo/dipra"
	"github.com/didikprabowo/dipra/middleware"
)

type (
	Personal struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	}
)

func main() {

	route := dipra.Default()
	route.Use(middleware.Logger())
	route.GET("/", func(c *dipra.Context) error {

		return c.JSON(http.StatusOK, dipra.M{
			"status": true,
		})
	})
	route.Run(":9020")
}
