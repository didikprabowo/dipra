package main

import (
	"github.com/didikprabowo/dipra"
)

func main() {
	r := dipra.Default()
	r.GET("/user/:id/:name", func(c *dipra.Context) error {
		id := c.Param("id")
		name := c.Param("name")
		return c.JSON(200, dipra.M{
			"ID":   id,
			"NAME": name,
		})
	})
	err := r.Run(":9020")
	if err != nil {
		panic(err)
	}
}
