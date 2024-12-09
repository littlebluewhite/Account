package convert

import "reflect"

func CreateFieldConvert(srcVal, dstVal reflect.Value) {
	if dstVal.IsValid() && srcVal.Type().AssignableTo(dstVal.Type()) && dstVal.CanSet() {
		dstVal.Set(srcVal)
		return
	}

	switch srcVal.Kind() {
	case reflect.Ptr:
		if srcVal.IsNil() {
			return
		}
		if dstVal.IsNil() {
			newDst := reflect.New(dstVal.Type().Elem())
			dstVal.Set(newDst)
		}
		CreateFieldConvert(srcVal.Elem(), dstVal.Elem())
	case reflect.Struct:
		for i := 0; i < srcVal.NumField(); i++ {
			field := srcVal.Type().Field(i)
			if field.PkgPath != "" {
				// Skip unexported fields
				continue
			}
			srcField := srcVal.Field(i)
			dstField := dstVal.FieldByName(field.Name)
			if dstField.IsValid() && dstField.CanSet() {
				CreateFieldConvert(srcField, dstField)
			}
		}
	case reflect.Slice, reflect.Array:
		if srcVal.IsNil() {
			return
		}
		if dstVal.Kind() != reflect.Slice && dstVal.Kind() != reflect.Array {
			return
		}
		newSlice := reflect.MakeSlice(dstVal.Type(), srcVal.Len(), srcVal.Cap())
		for i := 0; i < srcVal.Len(); i++ {
			CreateFieldConvert(srcVal.Index(i), newSlice.Index(i))
		}
		dstVal.Set(newSlice)
	default:
		// For basic types and unsupported kinds, attempt direct assignment
		if dstVal.IsValid() && srcVal.Type().AssignableTo(dstVal.Type()) && dstVal.CanSet() {
			dstVal.Set(srcVal)
		}
	}
}
