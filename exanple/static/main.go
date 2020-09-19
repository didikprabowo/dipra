package main

import (
	"github.com/didikprabowo/dipra"
)

func main() {
	r := dipra.Default()
	r.GET("/pps", func(c *dipra.Context) error {
		return c.File("public/p.png")
	})
	r.Static("/static", "./public")
	r.Run(":9000")
}
