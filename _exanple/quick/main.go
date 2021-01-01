package main

import (
	"fmt"

	"github.com/didikprabowo/dipra"
)

func main() {
	d := dipra.Default()

	d.GET("/quick", func(c *dipra.Context) error {
		return c.JSON(200, "WELCOME")
	})

	if err := d.Run(":9020"); err != nil {
		fmt.Println(err)
	}
}
