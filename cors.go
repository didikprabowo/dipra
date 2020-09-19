package dipra

import (
	"net/http"
	"strconv"
	"strings"
	"time"
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
func CORS() MiddlewareFunc {
	return BuildCors(DefaultConfig())
}

// CorsWithConifg is used setting cors
func CorsWithConifg(l CORSConfig) MiddlewareFunc {
	return BuildCors(l)
}

// BuildCors for used build cors
func BuildCors(cors CORSConfig) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) (err error) {
			methods := cors.SetCorsMethod()
			origings := cors.SetCorsOrigin()
			headers := cors.SetCorsHeader()
			age := cors.SetCorsMaxAge()
			credential := cors.SetCorsCredential()
			expose := cors.SetCorsExpose()

			res := c.GetResponse()
			res.Header().Add(string(HeaderVary), string(AccessControllReqHeaders))
			res.Header().Add(string(HeaderVary), string(AccessControllReqMethod))
			res.Header().Set(string(ACcessControllHeaders), headers)
			res.Header().Set(string(ACcessControllMethod), methods)
			res.Header().Set(string(AccessControllOrigin), origings)
			res.Header().Set(string(AccessControllMaxAge), age)
			res.Header().Set(string(AccessControllCredential), credential)
			res.Header().Set(string(AccessControllExposeHeaders), expose)
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
