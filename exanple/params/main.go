package main

import (
	"fmt"

	"github.com/didikprabowo/dipra"
)

func main() {
	r := dipra.Default()
	r.GET("/user/:id/:id", func(c *dipra.Context) error {
		id := c.Param("id")
		name := c.Param("name")
		fmt.Println("Hello", id, name)
		return nil
	})
	err := r.Run(":9020")
	if err != nil {
		panic(err)
	}
}
