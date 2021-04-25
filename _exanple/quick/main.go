package main

import (
	"fmt"

	"github.com/didikprabowo/dipra"
)

func main() {

	d := dipra.Default()

	d.GET("/didik/b", func(c *dipra.Context) error {
		return c.JSON(200, "hai "+c.URL.Path)
	})
	d.GET("/didik/a", func(c *dipra.Context) error {
		return c.JSON(200, "hai "+c.URL.Path)
	})
	// d.GET("/p", func(c *dipra.Context) error {
	// 	return c.JSON(200, "hai "+c.URL.Path)
	// })

	if err := d.Run(":9020"); err != nil {
		fmt.Println(err)
	}

}
