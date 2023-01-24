package emptiness

import (
	"reflect"
)

type IsZeroer interface {
	IsZero() bool
}

//nolint:cyclop
func IsZero(a interface{}) bool {
	v := reflect.ValueOf(a)
	rt := v.Type()
	for i := 0; i < rt.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
			if v.Field(i).Len() != 0 {
				return false
			}
		case reflect.Bool:
			if v.Field(i).Bool() {
				return false
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v.Field(i).Int() != 0 {
				return false
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			if v.Field(i).Uint() != 0 {
				return false
			}
		case reflect.Float32, reflect.Float64:
			if v.Field(i).Float() != 0 {
				return false
			}
		case reflect.Interface, reflect.Ptr:
			if !v.Field(i).IsNil() {
				return false
			}
		case reflect.Struct:
			if z, ok := v.Field(i).Interface().(IsZeroer); ok {
				if !z.IsZero() {
					return false
				}
			}
		}
	}
	return true
}
