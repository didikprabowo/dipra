package main

import (
	"github.com/didikprabowo/dipra"
)

func main() {
	r := dipra.Default()
	r.GET("/user/one/:id", func(c *dipra.Context) error {
		id := c.Param("id")
		return c.JSON(200, dipra.M{
			"ID": id,
		})
	})
	r.GET("/user/two/:id", func(c *dipra.Context) error {
		id := c.Param("id")
		return c.JSON(200, dipra.M{
			"ID": id,
		})
	})
	err := r.Run(":9020")
	if err != nil {
		panic(err)
	}
}
