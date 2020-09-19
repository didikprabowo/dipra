package dipra

import (
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
func (e *Node) ReserverURI() (out string) {
	var (
		newURI []string
	)

	rt := strings.Split(e.CleanPath(e.Route.Path), "/")
	uri := strings.Split(e.CleanPath(e.Ctx.RequestURI), "/")

	if len(rt) != len(uri) {
		return ""
	}

	for i := 0; i < len(rt); i++ {
		if strings.HasPrefix(rt[i], ":") {
			param := Param{
				Key:   rt[i][1:len(rt[i])],
				Value: uri[i],
			}
			rt[i] = uri[i]
			e.Ctx.Params.SetParam(param)
		}

		newURI = append(newURI, rt[i])
	}
	e.Fullpatch = newURI

	return strings.Join(newURI, "/")
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
