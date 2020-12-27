# Dipra mini framework golang


[![Build Status](https://img.shields.io/travis/com/didikprabowo/dipra/master?label=Build&logo=travis)](https://travis-ci.com/github/didikprabowo/dipra)
[![codecov](https://img.shields.io/codecov/c/github/didikprabowo/dipra/master?color=s&label=Codecov&logo=Codecov&token=27b8cb42a538455b8a64351bfb90010b)](https://codecov.io/gh/didikprabowo/dipra)
[![go-version](https://img.shields.io/github/go-mod/go-version/didikprabowo/dipra?logo=go)](https://github.com/didikprabowo/dipra)
[![release](https://img.shields.io/github/v/release/didikprabowo/dipra?label=Release&logo=Release)](https://github.com/didikprabowo/dipra/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/didikprabowo/dipra.svg)](https://pkg.go.dev/github.com/didikprabowo/dipra)


Welcome to dipra. Dipra is mini framework golang to build service application. Dipra have high performance speed. Makes native source codes. Suitable for build REST API.
        
## Feature
 - HTTP with all method and static file.
 - Data binding request body(JSON,XML) and query raw.
 - Response String,JSON,JSONP,XML and File.
 - Suport middlere handler for al and by route.
 - Enable and disable log.

## Documentation.
 - [Installation](#Installation)
 - [Routing](#routing)
 - [Grouping routes](#Grouping-routes)
 - [Hello world](#Hello-world)
 - [Request](#Request)
   - [Parameter Path](#Parameter-Path)
   - [Parameter Body](#Parameter-Bind-Body-Raw)
   - [Parameter Query](#Parameter-Query-string)
 - [Response](#Response)
 - [Static](#Static)
 - [Middleware](#Middleware)
   - [Logger](#logger)
   - [Cors](#cors) 
   - [Create Middleware](#Create-Middleware)


### Installation
Install must be have GO system on your PC. If you have it, you can install with cmd.
1. Install package
``` bash
go get -u github.com/didikprabowo/dipra
```

2. Import code
Kindly import package at top code program, for exampe : 
```go
import (
        "github.com/didikprabowo/dipra"
)

```
### Hello world

`main.go`
```go
package main

import (
	"net/http"
	"github.com/didikprabowo/dipra"
)

func main() {

	r := dipra.Default()

	r.GET("/hello-world", func(c *dipra.Context) error {
		return c.JSON(http.StatusOK, dipra.M{
			"data": "Hello world.",
		})
	})

	r.Run(":6000")
}
```

### Routing

Route base http server default golang. It leverages `sync pool` to use memory. By besides, dipra router will find priority routing which request it.

For use route must be define path and handler, but can be use middleware. For example : 


```go

r := dipra.Default()

r.GET("/GET", func(c *dipra.Context) error {
        return c.JSON(http.StatusOK, dipra.M{
                "data": "Hello Get.",
        })
})

r.POST("/POST", func(c *dipra.Context) error {
        return c.JSON(http.StatusOK, dipra.M{
                "data": "Hello Post.",
        })
})
```
### Grouping routes

```go
package main

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/didikprabowo/dipra"
)

func main() {

	r := dipra.Default()
	r.Use(dipra.Logger())

	// Normal Group
	v1 := r.Group("/v1")

	{
		// {basepath}/v1
		v1.GET("/", func(c *dipra.Context) error {
			return c.String(200, fmt.Sprintf("Welcome to api version v1 %+f\n", rand.Float64()))
		})

	}

	// Group use middleware
	v2 := r.Group("/v2", func(hf dipra.HandlerFunc) dipra.HandlerFunc {
		return func(c *dipra.Context) (err error) {
			if c.Query("name") == "jhon" {
				return c.String(http.StatusUnprocessableEntity, "Can't continue")
			}
			hf(c)
			return
		}
	})

	{
		// {basepath}/v2
		v2.GET("/", func(c *dipra.Context) error {
			return c.String(200, fmt.Sprintf("Welcome to api version v2 %+f\n", rand.Float64()))
		})

	}

	err := r.Run(":9020")
	if err != nil {
		panic(err)
	}
}



```
### Request
To get request from body raw you can bind with type mime(JSON,XML). By besides, you be able to bind query raw with format `?key=value&key=value`. For example used : 

#### Parameter Path
```go
// Example : /ping/1
r := dipra.Default()
r.GET("/ping/:id", func(c *dipra.Context) error {
        id := c.Param("id")
        return c.JSON(http.StatusOK, dipra.M{
                "data": id,
        })
})

```

#### Parameter Bind Body Raw

```go
r := dipra.Default()

r.GET("/json", func(c *dipra.Context) error {
        var p Personal
        err := c.ShouldJSON(&p)
        if err != nil {
                return c.JSON(http.StatusBadRequest, dipra.M{
                        "error": err.Error(),
                })
        }

        return c.JSON(http.StatusOK, dipra.M{
                "data": p,
        })
})
```
#### Parameter Query string

```go
// example : ?id=2&name=didik
r.GET("/query", func(c *dipra.Context) error {
        var p Personal
        singleID := p.Query("id")
        err := c.ShouldQuery(&p)
        if err != nil {
                return c.JSON(http.StatusBadRequest, dipra.M{
                        "error": err.Error(),
                })
        }
        return c.JSON(http.StatusOK, dipra.M{
                "data": p,
                "id" : singleID,
        })
})

```
### Response

Dipra can write response in the from of String,JSON,JSONP, XML and File. Fro example : 

```go 
r := dipra.Default()

r.GET("/json", func(c *dipra.Context) error {
        p := Personal{
                Name:    "Didik Prabowo",
                Address: "Wonogiri",
        }

        return c.JSON(http.StatusOK, dipra.M{
                "data": p,
        })
})

r.GET("/jsonp", func(c *dipra.Context) error {
        p := Personal{
                Name:    "Didik Prabowo",
                Address: "Wonogiri",
        }

        return c.JSONP(http.StatusOK, dipra.M{
                "data": p,
        })
})

r.GET("/string", func(c *dipra.Context) error {
        return c.String(200, "Welcome to dipra")
})

r.GET("/xml", func(c *dipra.Context) error {
        p := Personal{
                Name:    "Didik Prabowo",
                Address: "Wonogiri",
        }

        return c.XML(http.StatusOK, dipra.M{
                "data": p,
        })
})

r.GET("/file", func(c *dipra.Context) error {
        return c.File("public/p.png")
})

```

### Static

Dipra provide for you want get static file as .css, .html,.png and etc. For example : 

`main.go`
```go
package main

import (
	"github.com/didikprabowo/dipra"
)

func main() {
        r := dipra.Default()
          // {{basepath}}/get-image
	r.GET("/get-image", func(c *dipra.Context) error {
		return c.File("public/p.png")
        })
      
         // {{basepath}}/static/p.png
        r.Static("/static", "./public")
       
	r.Run(":9000")
}
```

### Middleware
Dipra provide middleware handle for handler function. For example you can print log in middleware. Any 2 method how to define middleware : Define in all route and spesific route.

#### Logger
```go
route := dipra.Default()
route.Use(dipra.Logger())
route.GET("/", func(c *dipra.Context) error {

        return c.JSON(http.StatusOK, dipra.M{
                "status": true,
        })
})
```

#### Cors

```go

route := dipra.Default()
// default :: route.Use(dipra.CORS())
route.Use(dipra.CorsWithConifg(dipra.CORSConfig{
        AllowOrigins: []string{"https://www.google.com"},
        AllowMethod:  []string{"*"},
        AllowHeaders: []string{"*"},
}))
route.GET("/", func(c *dipra.Context) error {

        return c.JSON(http.StatusOK, dipra.M{
                "data": "Example cors",
        })
})
```

#### Create Middleware

You can create new custom middleare at route or at all route.

```go

package main

import (
	"log"
	dipra "github.com/didikprabowo/dipra"
)

func main() {
	s := dipra.Default()
	s.Use(TesMiddlewareAll)
	s.GET("/ping", Get, TesMiddleware)
	err := s.Run(":6000")

	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}

}

func Get(c *dipra.Context) error {
	return c.JSON(200, dipra.M{
		"data": "tes",
	})
}

func TesMiddlewareAll(next dipra.HandlerFunc) dipra.HandlerFunc {
	return func(c *dipra.Context) error {
		return next(c)
	}
}

func TesMiddleware(next dipra.HandlerFunc) dipra.HandlerFunc {
	return func(c *dipra.Context) error {
		return next(c)
	}
}

```
## Credits
### Author 
- Didik Prabowo - https://kodingin.com

