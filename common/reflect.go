package common

import (
	"fmt"
	"reflect"
	"strconv"
)

// SetFieldByName sets the field
func SetFieldByName(s interface{}, field string, value string) (bool, error) {
	v := reflect.ValueOf(s).Elem()
	if !v.CanAddr() {
		return false, fmt.Errorf("cannot assign to the item passed, item must be a pointer in order to assign")
	}

	f := v.FieldByName(field)

	if !f.CanSet() {
		return false, nil
	}

	switch f.Kind() {
	case reflect.String:
		f.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return false, fmt.Errorf("invalid type of field %s, expecting type %s, err at conversion: %w", field, f.Type(), err)
		}

		f.SetInt(i)
	default:
		return false, fmt.Errorf("set field invalid for type %s", f.Type())
	}

	return true, nil
}
