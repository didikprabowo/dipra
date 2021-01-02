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
	route.Use(middleware.CorsWithConifg(middleware.CORSConfig{
		AllowOrigins: []string{"https://www.google.com"},
		AllowMethod:  []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	route.GET("/", func(c *dipra.Context) error {

		data := []Personal{
			Personal{
				Name:    "Didik",
				Address: "Wonogiri",
			},
			Personal{
				Name:    "Praboeo",
				Address: "Solo",
			},
		}

		return c.JSON(http.StatusOK, dipra.M{
			"data": data,
			"status": dipra.M{
				"message": "Berhasil",
			},
		})
	})
	route.Run(":9020")
}
