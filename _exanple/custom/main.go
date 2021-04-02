package main

import (
	"fmt"

	"github.com/didikprabowo/dipra"
)

func main() {

	route := dipra.Default()
	// route.GET("/", func(c *dipra.Context) error {
	// 	return c.JSON(200, "DIDIK")
	// 	// return nil
	// })
	// route.Any("/", func(c *dipra.Context) error {
	// 	// return c.Ha
	// 	c.Han
	// 	// return nil
	// })
	s := route.Group("/didik")
	s.GET("/tes", func(c *dipra.Context) (err error) {
		fmt.Println(c.GetHeader("*"))
		return fmt.Errorf("tes")
	})
	// route.GET("/didik", func(c *dipra.Context) (err error) {
	// 	return c.JSON(300, "OKE")
	// })

	route.Run(":9020")
}
