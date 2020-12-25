package main

import (
	"github.com/didikprabowo/dipra"
)

func main() {
	r := dipra.Default()
	r.Use(dipra.Logger())

	// Example image
	r.GET("/image", func(c *dipra.Context) error {
		return c.File("public/p.png")
	})

	// Example Doc
	r.GET("/doc", func(c *dipra.Context) error {
		return c.File("public/index.html")
	})

	r.Static("/static", "./public")
	r.Run(":9020")
}
