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
				path:       n.path[i:],
				label:      n.label,
				nodeType:   static,
				children:   n.children,
				HandleFunc: n.HandleFunc,
			}
			n.label = string(n.path[i])
			n.children = []*node{child}
			n.path = path[:i]
			n.HandleFunc = nil
		}

		if i < ls {
			path = path[i:]

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

func (n *node) find(path string) (ok bool, fpath string, HandleFunc HandlerFunc) {
	if path == "" {
		panic("dipra : path is empty")
	}

walk:
	for {
		search := n.path

		if len(path) > len(search) {
			if path[:len(n.path)] == (n.path) {

				path = path[len(search):]
				fpath += search
				prefix := path[0]

				if !n.isWildChild {
					for i, v := range n.label {
						if byte(v) == prefix {
							n = n.children[i]
							continue walk
						}
					}
					return
				}
			}
		} else if search == path {
			HandleFunc = n.HandleFunc
			ok = true
			fpath += search
		}

		return
	}
}
