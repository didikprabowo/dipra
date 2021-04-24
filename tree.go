package dipra

type (
	nodeType uint8
	node     struct {
		nodeType    nodeType
		path        string
		label       string
		children    []*node
		HandleFunc  HandlerFunc
		isWildChild bool
	}
)

const (
	static nodeType = iota
	root
	param
)

func (n *node) insert(method, path string, h HandlerFunc) {
	if path[0] != '/' {
		path = "/" + path
	}

	if n == nil {
		panic("dipra : invalid node")
	}

	if h == nil {
		panic("dipra : Hander must be not nil")
	}

	if n.path == "" && n.label == "" {
		n.insertChild(path, h)
		n.nodeType = root
		return
	}

walk:
	for {
		ls := len(path)
		lp := len(n.path)

		max := lp
		if ls < max {
			max = ls
		}

		i := 0

		for ; i < max && n.path[i] == path[i]; i++ {
		}

		if i < lp {

			child := &node{
				path:        n.path[i:],
				label:       n.label,
				nodeType:    static,
				children:    n.children,
				HandleFunc:  n.HandleFunc,
				isWildChild: n.isWildChild,
			}
			n.label = string(n.path[i])
			n.children = []*node{child}
			n.path = path[:i]
			n.HandleFunc = nil
		}

		if i < ls {

			path = path[i:]

			if n.isWildChild {
				if len(path) >= len(n.path) &&
					n.path == path[:len(n.path)] &&
					(len(n.path) >= len(path) || path[len(n.path)] == '/') {
					continue walk
				} else {
					panic("path is conflict")
				}
			}

			prefix := path[0]
			for _, v := range n.label {
				if byte(v) == prefix {

					cn := n.findChildren(byte(v))
					if cn != nil {
						n = cn
						continue walk
					}
				}
			}

			if prefix != ':' {
				n.label += string(prefix)
				child := &node{}
				n.children = append(n.children, child)
				n = child
			}
			n.insertChild(path, h)

			return
		}

		return
	}

}

func (n *node) insertChild(path string, h HandlerFunc) {
	for {
		i, wildcard := n.findWildcard(path)

		if i < 0 {
			break
		}

		if len(n.children) > 0 {
			panic("dipra : path is conflict")
		}

		if wildcard[0] == ':' {
			n.path = path[:i]
			path = path[i:]

			child := &node{
				nodeType: param,
				path:     wildcard,
			}
			n.isWildChild = true

			n.children = []*node{child}
			n = child

			if len(wildcard) < len(path) {
				path = path[len(wildcard):]
				child := &node{}
				n.children = []*node{child}
				n = child
				continue
			}
			n.insertHandle(h)
			return
		}
		return
	}

	n.insertHandle(h)
	n.insertPath(path)
}

func (n *node) findWildcard(path string) (i int, wilchard string) {

	for i, l := 0, len(path); i < l; i++ {
		if path[i] == ':' {
			j := i
			for ; i < l && path[i] != '/'; i++ {
			}

			return j, path[j:i]
		}
	}

	return -1, ""
}

func (n *node) findChildren(prefix byte) *node {
	for _, v := range n.children {
		if v.path[0] == prefix {
			return v
		}
	}
	return nil
}

func (n *node) insertPath(path string) {
	n.path = path
}

func (n *node) insertHandle(h HandlerFunc) {
	n.HandleFunc = h
}

func (n *node) find(path string, params *params) (ok bool, ps *params, HandleFunc HandlerFunc) {
	if path == "" {
		panic("dipra : path is empty")
	}

walk:
	for {
		// search := path
		cpath := n.path
		if len(path) > len(cpath) {

			if path[:len(cpath)] == cpath {

				path = path[len(cpath):]

				prefix := path[0]

				if !n.isWildChild {
					for i, v := range n.label {
						if byte(v) == prefix {
							n = n.children[i]
							continue walk
						}
					}

					// check route handler

					if path == "/" && n.HandleFunc != nil {
						ok = true
					}
					return
				}

				n = n.children[0]
				switch n.nodeType {
				case param:
					i := 0
					lp := len(path)
					for ; i < lp && path[i] != '/'; i++ {
					}

					if params != nil {
						if ps == nil {
							ps = params
						}

						*ps = append(*ps, viewParam{
							Key:   n.path[1:],
							Value: path[:i],
						})
					}

					if i < len(path) {
						if len(n.children) > 0 {
							path = path[i:]
							n = n.children[0]
							continue walk
						}

						ok = (len(path) == i+1)
						return
					}

					if n.HandleFunc != nil {
						HandleFunc = n.HandleFunc
						return
					} else if len(n.children) == 1 {
						n = n.children[0]
						ok = (n.path == "/" && n.HandleFunc != nil) || (n.path == "" && n.label == "/")
					}

					return
				}
			}
		} else if cpath == path {
			HandleFunc = n.HandleFunc
			ok = true
		}

		return
	}
}
