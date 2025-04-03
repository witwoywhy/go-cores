package utils

import "reflect"

var (
	_true interface{} = true
)

func MakeMap[K comparable, V any, T any, E ~[]T](items E, keyField, valueField string) map[K]V {
	response := make(map[K]V)

	for _, v := range items {
		elem := reflect.ValueOf(v)

		k := elem.FieldByName(keyField)

		if !k.IsValid() {
			continue
		}

		key, ok := k.Interface().(K)
		if !ok {
			continue
		}

		switch valueField {
		case "true":
			response[key] = _true.(V)
		case "struct":
			response[key] = elem.Interface().(V)
		default:
			v := elem.FieldByName(valueField)
			if !v.IsValid() {
				continue
			}

			val, ok := v.Interface().(V)
			if ok {
				response[key] = val
			}
		}
	}

	return response
}
