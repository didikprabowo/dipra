package main

import (
	"net/http"

	"github.com/didikprabowo/dipra"
)

type (
	Personal struct {
		Name    string `json:"name" xml:"name"`
		Address string `json:"address" xml:"address"`
	}
)

func main() {

	r := dipra.Default()

	r.GET("/json", func(c *dipra.Context) error {
		p := Personal{
			Name:    "Didik Prabowo",
			Address: "Wonogiri",
		}

		return c.JSON(http.StatusOK, dipra.M{
			"data": p,
		})
	})

	r.GET("/jsonp", func(c *dipra.Context) error {
		p := Personal{
			Name:    "Didik Prabowo",
			Address: "Wonogiri",
		}

		return c.JSONP(http.StatusOK, dipra.M{
			"data": p,
		})
	})

	r.GET("/string", func(c *dipra.Context) error {
		return c.String(200, "Welcome to dipra")
	})

	r.GET("/xml", func(c *dipra.Context) error {
		p := Personal{
			Name:    "Didik Prabowo",
			Address: "Wonogiri",
		}

		return c.XML(http.StatusOK, dipra.M{
			"data": p,
		})
	})

	err := r.Run(":7000")
	if err != nil {
		panic(err)
	}
}
