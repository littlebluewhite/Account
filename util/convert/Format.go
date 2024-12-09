package convert

import "reflect"

func Format[T, U any](m []T) []U {
	result := make([]U, len(m))
	for _, item := range m {
		i := Model2Entry[T, U](item)
		result = append(result, i)
	}
	return result
}

func Model2Entry[T, U any](m T) U {
	var result U
	srcVal := reflect.ValueOf(&m).Elem()
	dstVal := reflect.ValueOf(&result).Elem()
	CreateFieldConvert(srcVal, dstVal)
	return result
}
