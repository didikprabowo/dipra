package main

import (
	"net/http"

	"github.com/didikprabowo/dipra"
)

type (
	Personal struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	}
)

func main() {

	route := dipra.Default()
	route.Use(dipra.Logger())
	route.GET("/", func(c *dipra.Context) error {

		return c.JSON(http.StatusOK, dipra.M{
			"status": true,
		})
	})
	route.Run(":9000")
}
