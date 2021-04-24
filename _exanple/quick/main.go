package main

import (
	"log"
	"net/http"

	"github.com/didikprabowo/dipra"
)

func main() {

	d := dipra.Default()
	// d.Use(middleware.Logger())
	// d.Use(func(hf dipra.HandlerFunc) dipra.HandlerFunc {
	// 	return func(c *dipra.Context) error {
	// 		fmt.Fprintf(os.Stdout, fmt.Sprintf("Log %s - %s", c.Request.URL.Path, c.Param("didik")))
	// 		return hf(c)
	// 	}
	// })

	d.POST("/didik", func(c *dipra.Context) error {
		return c.JSON(200, c.Param("didik")+"s")
	}, func(hf dipra.HandlerFunc) dipra.HandlerFunc {
		return func(c *dipra.Context) error {
			return hf(c)
		}
	})

	srv := &http.Server{
		Addr:    ":9020",
		Handler: d,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Printf("listen: %s\n", err)
	}

	// if err := d.Run(":9020"); err != nil {
	// 	fmt.Println(err)
	// }

}
