package main

import (
	"github.com/didikprabowo/dipra"
)

func main() {
	r := dipra.Default()
	r.GET("/didikprabowo", func(c *dipra.Context) (err error) {
		return c.File("public/p.pngs")
		return err
	})
	r.GET("/p", func(c *dipra.Context) (err error) {
		return c.JSON(200, "s")
	})
	// r.Static("/static", "./public")
	r.Run(":9020")
}
