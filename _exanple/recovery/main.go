package main

import (
	"github.com/didikprabowo/dipra"
	"github.com/didikprabowo/dipra/middleware"
)

func main() {
	r := dipra.Default()
	r.Use(middleware.Recovery())
	r.GET("/", func(c *dipra.Context) error {
		panic("Error panic")
	})

	err := r.Run(":9020")
	if err != nil {
		panic(err)
	}
}
