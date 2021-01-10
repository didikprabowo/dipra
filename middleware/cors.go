package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/didikprabowo/dipra"
)

type (
	// CORSConfig ...
	CORSConfig struct {
		// AllowMethod
		AllowMethod []string
		// AllowOrigins
		AllowOrigins []string
		// AllowHeaders
		AllowHeaders []string
		// Max age
		MaxAge time.Duration
		// AllowCredential
		AllowCredential bool
		// ExposeHeader
		ExposeHeader []string
	}
)

// DefaultConfig Cors
func DefaultConfig() CORSConfig {
	return CORSConfig{
		AllowMethod: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodDelete,
			http.MethodPatch,
			http.MethodPut,
		},
		AllowOrigins: []string{
			"*",
		},
		AllowCredential: false,
		MaxAge:          12 * time.Hour,
	}
}

// CORS use Cors
func CORS() dipra.MiddlewareFunc {
	return BuildCors(DefaultConfig())
}

// CorsWithConifg is used setting cors
func CorsWithConifg(l CORSConfig) dipra.MiddlewareFunc {
	return BuildCors(l)
}

// BuildCors for used build cors
func BuildCors(cors CORSConfig) dipra.MiddlewareFunc {
	return func(h dipra.HandlerFunc) dipra.HandlerFunc {
		return func(c *dipra.Context) (err error) {
			methods := cors.SetCorsMethod()
			origings := cors.SetCorsOrigin()
			headers := cors.SetCorsHeader()
			age := cors.SetCorsMaxAge()
			credential := cors.SetCorsCredential()
			expose := cors.SetCorsExpose()

			res := c.GetResponse()
			res.Header().Add(string(dipra.HeaderVary), string(dipra.AccessControllReqHeaders))
			res.Header().Add(string(dipra.HeaderVary), string(dipra.AccessControllReqMethod))
			res.Header().Set(string(dipra.ACcessControllHeaders), headers)
			res.Header().Set(string(dipra.ACcessControllMethod), methods)
			res.Header().Set(string(dipra.AccessControllOrigin), origings)
			res.Header().Set(string(dipra.AccessControllMaxAge), age)
			res.Header().Set(string(dipra.AccessControllCredential), credential)
			res.Header().Set(string(dipra.AccessControllExposeHeaders), expose)
			h(c)
			return err
		}
	}
}

// SetCorsMethod for set method
func (crs *CORSConfig) SetCorsMethod() string {
	return strings.Join(crs.AllowMethod, ",")
}

// SetCorsHeader for set header
func (crs *CORSConfig) SetCorsHeader() string {
	return strings.Join(crs.AllowHeaders, ",")
}

// SetCorsOrigin for set origin allow
func (crs *CORSConfig) SetCorsOrigin() string {
	return strings.Join(crs.AllowOrigins, ",")
}

// SetCorsExpose for expose header
func (crs *CORSConfig) SetCorsExpose() string {
	return strings.Join(crs.ExposeHeader, ",")
}

// SetCorsMaxAge for set max age
func (crs *CORSConfig) SetCorsMaxAge() string {
	age := strconv.Itoa(int(crs.MaxAge))
	if len(age) == 0 {
		return "3600"
	}

	return age
}

// SetCorsCredential set allow creadential
func (crs *CORSConfig) SetCorsCredential() string {
	return strconv.FormatBool(crs.AllowCredential)
}
