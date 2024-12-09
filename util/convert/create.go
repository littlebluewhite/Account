package convert

import (
	"reflect"
)

func CreateConvert[T, U any](c []*U) []*T {
	result := make([]*T, 0, len(c))
	for _, item := range c {
		srcVal := reflect.ValueOf(item).Elem()
		var resultItem T
		dstVal := reflect.ValueOf(&resultItem).Elem()
		CreateFieldConvert(srcVal, dstVal)
		result = append(result, &resultItem)
	}
	return result
}
