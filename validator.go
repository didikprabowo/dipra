package dipra

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var (
	// _v ...
	_v Validate = Validator{}
	// NumericRgx ...
	NumericRgx = regexp.MustCompile("^[0-9]+$")
	// TextRgx ...
	TextRgx = regexp.MustCompile("^[0-9a-zA-Z]+$")
	// EmailRgx ...
	EmailRgx = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	// UUID4Rgx ...
	UUID4Rgx = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	// UUIDRgx ...
	UUIDRgx = regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$")
	// UppercaseRgx ...
	UppercaseRgx = regexp.MustCompile("^[A-Z]+$")
	// LowercaseRgx ...
	LowercaseRgx = regexp.MustCompile("^[a-z]+$")
)

const (
	// ValidateRequired ...
	ValidateRequired string = "required"
	// TagName ...
	TagName string = "is_valid"
	// TypeNotFound ...
	TypeNotFound string = "Type not found, Available for struct"
)

type (
	// Validate method
	Validate interface {

		// Validate provide handler all validator
		Validate(v interface{}) error

		// IsEmpty use check data is empty
		IsEmpty(data string) (is bool)

		// IsNumeric use check data is numeric
		IsNumeric(data string) (is bool)

		// IsEmail use check data is email
		IsEmail(data string) (is bool)

		// UUID4 use check data type is uuid4
		IsUUID4(data string) (is bool)

		// IsUpperCase use check data string is upper
		IsUppercase(data string) (is bool)

		// IsLowerCase use check data string is lower
		IsLowerCase(data string) (is bool)
	}
	// Validator struct
	Validator struct {
		Type      reflect.Type
		ChanValue interface{}
		Error     ValidatorMessage
	}
	// ValidatorMessage for handler error intance
	ValidatorMessage struct {
		Code     int
		Validate []string
		Internal string
	}
)

// HandlerValid
// type HandlerValid func(string) bool

// TypeValid ...
var TypeValid = map[string]func(string) bool{
	"email": _v.IsEmail,
}

func NewValidator() Validate {
	var v Validate = Validator{}
	return v
}

func (e ValidatorMessage) Error() (errMessage string) {

	if e.Internal != "" {
		errMessage += fmt.Sprintf("%v", e.Internal)
	}

	if len(e.Validate) > 0 {
		errMessage += fmt.Sprintf("%+v ", strings.Join(e.Validate, ","))
	}

	return errMessage
}

// Validate ...
func (v Validator) Validate(data interface{}) (err error) {
	value := reflect.ValueOf(data)

	if value.Kind() == reflect.Struct {
		err := v.ValidateStruct(value)
		if err != nil {
			return err
		}
	}

	return err
}

// ValidateStruct ...
func (v Validator) ValidateStruct(r reflect.Value) (err error) {
	if r.Kind() != reflect.Struct {
		return ValidatorMessage{Internal: TypeNotFound,
			Code: 500}
	}
	for i := 0; i < r.NumField(); i++ {
		// typeV := r.Type().Field(i)
		// valueField := r.Field(i).String()
		// v.checkRequired(typeV.Tag.Get(TagName), valueField, typeV.Name)
	}

	if len(v.Error.Validate) > 0 || v.Error.Internal != "" {
		v.Error.Code = 500
		return v.Error
	}

	return err
}

func (v Validate) splitValue(tags reflect.StructTag) (s string) {
	tag := strings.Split(tags, ",")
	for i := 0; i < len(tags); i++ {
		s = tag[i]
	}
	return s
}

// IsValidate ...
func (v Validator) checkRequired(data string, vs string, name string) {
	datas := strings.Split(data, ",")
	for i := 0; i < len(datas); i++ {
		if datas[i] == ValidateRequired {
			if vs == "" {
				errs := fmt.Sprintf("%v : This fill is required", name)
				v.Error.Validate = append(v.Error.Validate, errs)
			}
		}
	}
}

// IsEmpty ...
func (v Validator) IsEmpty(data string) (is bool) {
	if data != "" {
		return true
	}
	return false
}

// IsNumeric ...
func (v Validator) IsNumeric(data string) (is bool) {
	if v.IsEmpty(data) {
		return NumericRgx.MatchString(data)
	}
	return is
}

// IsEmail ...
func (v Validator) IsEmail(data string) (is bool) {
	if v.IsEmpty(data) {
		return EmailRgx.MatchString(data)
	}
	return is
}

// IsUUI ...
func (v Validator) IsUUI(data string) (is bool) {
	if v.IsEmpty(data) {
		return UUIDRgx.MatchString(data)
	}
	return is
}

// UUID4 ...
func (v Validator) IsUUID4(data string) (is bool) {
	if v.IsEmpty(data) {
		return UUID4Rgx.MatchString(data)
	}
	return is
}

// IsText ...
func (v Validator) IsText(data string) (is bool) {
	if v.IsEmpty(data) {
		return TextRgx.MatchString(data)
	}
	return is
}

// IsUppercase ...
func (v Validator) IsUppercase(data string) (is bool) {
	if v.IsEmpty(data) {
		return UppercaseRgx.MatchString(data)
	}
	return is
}

// IsLowerCase ...
func (v Validator) IsLowerCase(data string) (is bool) {
	if v.IsEmpty(data) {
		return UppercaseRgx.MatchString(data)
	}
	return is
}

// validateType ...
func (v Validator) validateType(data string) bool {
	return TypeValid[string(v.Type)](data)
}
