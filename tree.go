package dipra

import (
	"fmt"
	"net/http"
	"path"
	"strings"
)

type (
	// Node ...
	Node struct {
		Ctx       *Context
		Route     Route
		Fullpatch []string
	}
)

// SetNode ...
func (e *Node) SetNode(ctx *Context, Route Route) {
	e.Ctx = ctx
	e.Route = Route
}

// ReserverURI ...
func (e *Node) ReserverURI() (out string, err error) {

	e.Ctx.Params.clean()

	var (
		newURI []string
	)

	rt := strings.Split(e.CleanPath(e.Route.Path), "/")
	uri := strings.Split(e.CleanPath(e.Ctx.RequestURI), "/")

	if len(rt) != len(uri) {
		return "", err
	}
	Qpath := make(map[string]bool)

	index := 0
	for i := 0; i < len(rt); i++ {

		if strings.HasPrefix(rt[i], ":") {
			param := Param{
				Key:   rt[i][1:len(rt[i])],
				Value: uri[i],
			}

			switch {
			case Qpath[param.Key]:
				defaulterrorHTTP(e.Ctx.ResponseWriter, http.StatusInternalServerError, fmt.Errorf("Can't use double parameter"))
				return out, fmt.Errorf("Can't use double parameter")
			default:
				rt[i] = uri[i]
				index++
				e.Ctx.Params.SetParam(param)
				Qpath[param.Key] = true
				newURI = append(newURI, rt[i])
			}
		}

	}
	e.Fullpatch = newURI

	return strings.Join(newURI, "/"), err
}

// CleanPath ...
func (e *Node) CleanPath(r string) string {
	r = path.Clean(r)
	r = strings.TrimPrefix(r, "/")
	r = strings.TrimSuffix(r, "/")
	if len(r) == 0 {
		return ""
	}
	return r
}
