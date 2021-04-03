package dipra

import (
	"errors"
	"net/http"
	"path"
	"regexp"
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
		uries []string
		rt    = strings.Split(e.CleanPath(e.Route.Path), "/")
		uri   = strings.Split(e.CleanPath(e.Ctx.RequestURI), "/")
		qpath = make(map[string]bool)
	)

	index := 0
	for i := 0; i < len(rt); i++ {
		if rt[i] == "" {
			continue
		}

		switch {
		case strings.HasPrefix(rt[i], ":"):

			param := Param{
				Key:   rt[i][1:len(rt[i])],
				Value: uri[i],
			}

			switch {
			case qpath[param.Key]:
				return out, errors.New("Can't used double parameter")
			default:
				rt[i] = uri[i]
				index++
				e.Ctx.Params.SetParam(param)
				qpath[param.Key] = true
			}
		case strings.HasPrefix(rt[i], "*"):

			rt[i] = strings.Join(uri[i:(len(uri))], "/")
		}

		if rt[i] == "" {
			err = errors.New(http.StatusText(http.StatusNotFound))
			return "", err
		}

		// Validation Request URL
		pathValidation := func(path string) (err error) {
			isAllow := regexp.MustCompile(`[a-zA-Z0-9]$`)
			ok := isAllow.MatchString(path)
			if !ok {
				err = errors.New("Character URL not allowed")
			}
			return err
		}

		err = pathValidation(rt[i])

		uries = append(uries, rt[i])
	}

	e.Fullpatch = uries

	path := strings.Join(uries, "/")
	pacher := e.Route.Path

	e.Ctx.SetPath("/" + path)
	e.Ctx.SetPatcher(pacher)

	return path, err
}

// CleanPath ...
func (e *Node) CleanPath(r string) string {
	r = path.Clean(r)
	r = strings.TrimPrefix(r, "/")
	r = strings.TrimSuffix(r, "/")
	r = strings.TrimSpace(r)
	if len(r) == 0 {
		return ""
	}
	return r
}
