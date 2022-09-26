package common

import (
	"fmt"
	"reflect"
	"strconv"
)

// SetFieldByName sets the field.
func SetFieldByName(s interface{}, key string, value string) (bool, error) {
	v := reflect.ValueOf(s).Elem()
	if !v.CanAddr() {
		return false, fmt.Errorf("cannot assign to the item passed, item must be a pointer in order to assign")
	}

	field := v.FieldByName(key)

	if !field.CanSet() {
		return false, nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return false, fmt.Errorf("invalid type of field %s, expecting type %s, err at conversion: %w", key, field.Type(), err)
		}

		field.SetInt(i)
	default:
		return false, fmt.Errorf("set field invalid for type %s", field.Type())
	}

	return true, nil
}
