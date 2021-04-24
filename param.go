package dipra

// Param ...
type (
	viewParam struct {
		Key   string
		Value string
	}

	params []viewParam
)

func (p *params) listParam() *params {
	return p
}

func (p *params) putParams(v *params) {

	if v == nil {
		return
	}

	*p = append(*p, *v...)
}

func (p *params) getParam(search string) string {
	for _, v := range *p {
		if v.Key == search {
			return v.Value
		}
	}

	return ""
}
