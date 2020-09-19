package main

import (
	"net/http"

	"github.com/didikprabowo/dipra"
)

type (
	Personal struct {
		Name    string `json:"name" xml:"name" yml:"name"`
		Address string `json:"address" xml:"address" yaml:"address"`
	}
)

func main() {

	r := dipra.Default()

	r.GET("/json", func(c *dipra.Context) error {
		var p Personal
		err := c.ShouldJSON(&p)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dipra.M{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, dipra.M{
			"data": p,
		})
	})

	r.GET("/xml", func(c *dipra.Context) error {
		var p Personal
		err := c.ShouldXML(&p)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dipra.M{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusOK, dipra.M{
			"data": p,
		})
	})

	r.GET("/query", func(c *dipra.Context) error {
		var p Personal
		err := c.ShouldQuery(&p)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dipra.M{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusOK, dipra.M{
			"data": p,
		})
	})

	r.Run(":6000")
}
