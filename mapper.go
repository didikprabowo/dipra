package dipra

import (
	"errors"
	"reflect"
)

type (
	// Field temporary struct
	Field struct {
		Key   string
		Raw   string
		Type  string
		Value interface{}
	}
	// Mapper for instance
	mapper struct {
		Values reflect.Value
		Type   reflect.Type
		Field  []Field
		Out    interface{}
	}
)

// Set Field key value
func (key *mapper) Set(v interface{}) (err error) {
	key.Out = v
	var valuein = reflect.ValueOf(v)
	if valuein.Kind() == reflect.Struct {
		return errors.New("Type not suport, please use type prt(&)")
	}
	key.Values = valuein.Elem()
	key.Type = key.Values.Type()
	for i := 0; i < key.Values.NumField(); i++ {
		key.Field = append(key.Field, Field{
			Key:   key.Type.Field(i).Name,
			Raw:   key.Type.Field(i).Tag.Get("json"),
			Type:  key.Type.Field(i).Type.String(),
			Value: key.Values.Field(i).Interface(),
		})
	}

	return err
}

// MapToStruct is used convert map to struct
func (key *mapper) MapToStruct(m map[string]string) (err error) {
	for k, vs := range m {
		for _, kp := range key.Field {
			if kp.Raw == k {
				name := reflect.ValueOf(key.Out).Elem().FieldByName(kp.Key)
				if name.IsValid() {
					if name.CanSet() {
						name.Set(reflect.ValueOf(vs))
					}
				} else {
					err = errors.New("Invalid properties")
				}
			}
		}
	}
	return err
}
