package main

import (
	"github.com/didikprabowo/dipra"
)

func main() {
	r := dipra.Default()
	r.Use(dipra.Recovery())
	r.GET("/", func(c *dipra.Context) error {
		panic("Error panic")
	})

	err := r.Run(":9020")
	if err != nil {
		panic(err)
	}
}
