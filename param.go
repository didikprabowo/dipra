package dipra

// Param ...
type Param struct {
	Key   string
	Value string
}

var params []Param

// SetParam ...
func (p *Param) SetParam(param Param) {
	params = append(params, param)
}

// Param ...
func (p *Param) Param(b string) string {
	return p.ParamByName(b)
}

// ParamByName ...
func (p *Param) ParamByName(b string) string {
	for _, v := range params {
		if v.Key == b {
			return v.Value
		}
	}
	return ""
}

// GetParam ...
func (p *Param) GetParam() *Param {
	return p
}
