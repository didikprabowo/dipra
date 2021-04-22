package main

import (
	"fmt"
	"os"

	"github.com/didikprabowo/dipra"
)

func main() {
	d := dipra.Default()

	d.GET("/", func(c *dipra.Context) error {
		return c.JSON(200, "WELCOME")
	}, func(hf dipra.HandlerFunc) dipra.HandlerFunc {
		return func(c *dipra.Context) error {
			fmt.Fprintf(os.Stdout, fmt.Sprintf("halo %s", c.Request.URL.Path))
			return hf(c)
		}
	})

	if err := d.Run(":9020"); err != nil {
		fmt.Println(err)
	}
}
