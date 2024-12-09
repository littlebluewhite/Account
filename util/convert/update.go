package convert

import (
	"fmt"
	"reflect"
)

type UpdateSchema[T comparable] interface {
	GetKey(key string) T
}

func UpdateConvert2[S comparable, T any, U UpdateSchema[S]](modelMap map[S]T, uct []U, key string) ([]*T, error) {
	result := make([]*T, 0, len(uct))
	for _, u := range uct {
		c, ok := modelMap[(u).GetKey(key)]
		if !ok {
			return nil, fmt.Errorf("key not found")
		}
		srcVal := reflect.ValueOf(u).Elem()
		dstVal := reflect.ValueOf(&c).Elem()
		UpdateFieldConvert(srcVal, dstVal)
		result = append(result, &c)
	}
	return result, nil
}

func UpdateFieldConvert(srcVal, dstVal reflect.Value) {
	// Direct assignment for assignable types
	if dstVal.IsValid() && srcVal.Type().AssignableTo(dstVal.Type()) && dstVal.CanSet() {
		dstVal.Set(srcVal)
		return
	}

	switch srcVal.Kind() {
	case reflect.Ptr:
		if srcVal.IsNil() {
			return
		}
		if dstVal.Kind() == reflect.Ptr {
			if dstVal.IsNil() {
				newDst := reflect.New(dstVal.Type().Elem())
				dstVal.Set(newDst)
			}
			UpdateFieldConvert(srcVal.Elem(), dstVal.Elem())
		} else {
			UpdateFieldConvert(srcVal.Elem(), dstVal)
		}
	case reflect.Struct:
		for i := 0; i < srcVal.NumField(); i++ {
			srcField := srcVal.Field(i)
			fieldName := srcVal.Type().Field(i).Name
			dstField := dstVal.FieldByName(fieldName)
			if dstField.IsValid() && dstField.CanSet() {
				UpdateFieldConvert(srcField, dstField)
			}
		}
	case reflect.Slice:
		if srcVal.IsNil() {
			return
		}
		newSlice := reflect.MakeSlice(dstVal.Type(), srcVal.Len(), srcVal.Cap())
		for i := 0; i < srcVal.Len(); i++ {
			UpdateFieldConvert(srcVal.Index(i), newSlice.Index(i))
		}
		dstVal.Set(newSlice)
	default:
		if dstVal.IsValid() && dstVal.Type() == srcVal.Type() && dstVal.CanSet() {
			dstVal.Set(srcVal)
		}
	}
}

func UpdateConvert[T any, U any](modelSlice []T, uct []U, key string) ([]T, error) {
	modelMap := make(map[interface{}]*T)
	for i := range modelSlice {
		m := &modelSlice[i]
		keyValue, err := getKeyValue(*m, key)
		if err != nil {
			return nil, err
		}
		modelMap[keyValue] = m
	}

	for _, u := range uct {
		uKeyValue, err := getKeyValue(u, key)
		if err != nil {
			return nil, err
		}

		if modelPtr, exists := modelMap[uKeyValue]; exists {
			err := updateModel(modelPtr, u, key)
			if err != nil {
				return nil, err
			}
		} else {
			newModel, err := convertUpdateToModel[T](u)
			if err != nil {
				return nil, err
			}
			modelMap[uKeyValue] = &newModel
		}
	}

	updatedModels := make([]T, 0, len(modelMap))
	for _, modelPtr := range modelMap {
		updatedModels = append(updatedModels, *modelPtr)
	}

	return updatedModels, nil
}

// getKeyValue retrieves the key value from an object.
func getKeyValue(obj interface{}, key string) (interface{}, error) {
	v := reflect.ValueOf(obj)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	field := v.FieldByName(key)
	if !field.IsValid() {
		return nil, fmt.Errorf("key '%s' not found in struct", key)
	}
	if !field.Type().Comparable() {
		return nil, fmt.Errorf("key '%s' is not comparable", key)
	}
	return field.Interface(), nil
}

// updateModel updates the fields of the model based on the update provided.
func updateModel(modelPtr interface{}, update interface{}, key string) error {
	modelValue := reflect.ValueOf(modelPtr)
	if modelValue.Kind() != reflect.Ptr {
		return fmt.Errorf("modelPtr must be a pointer")
	}
	modelValue = modelValue.Elem()

	updateValue := reflect.ValueOf(update)
	for updateValue.Kind() == reflect.Ptr || updateValue.Kind() == reflect.Interface {
		updateValue = updateValue.Elem()
	}

	if modelValue.Kind() != reflect.Struct || updateValue.Kind() != reflect.Struct {
		return fmt.Errorf("model and update must be structs")
	}

	updateType := updateValue.Type()
	for i := 0; i < updateType.NumField(); i++ {
		field := updateType.Field(i)
		if field.Name == key {
			continue // Skip key field
		}
		uFieldValue := updateValue.Field(i)
		mFieldValue := modelValue.FieldByName(field.Name)
		if !mFieldValue.IsValid() || !mFieldValue.CanSet() {
			continue
		}

		if isNilOrZero(uFieldValue) {
			continue // Skip nil or zero value fields
		}

		if mFieldValue.Kind() == reflect.Struct && uFieldValue.Kind() == reflect.Struct {
			if !mFieldValue.CanAddr() {
				return fmt.Errorf("cannot get address of field %s", field.Name)
			}
			err := updateModel(mFieldValue.Addr().Interface(), uFieldValue.Interface(), key)
			if err != nil {
				return err
			}
		} else if mFieldValue.Kind() == reflect.Slice && uFieldValue.Kind() == reflect.Slice {
			elemType := mFieldValue.Type().Elem()
			if elemType.Kind() == reflect.Struct {
				updatedSlice, err := updateSlice(mFieldValue.Interface(), uFieldValue.Interface(), key)
				if err != nil {
					return err
				}
				mFieldValue.Set(reflect.ValueOf(updatedSlice))
			} else {
				mFieldValue.Set(uFieldValue)
			}
		} else {
			// Updated else clause
			if mFieldValue.Type() == uFieldValue.Type() && mFieldValue.CanSet() {
				mFieldValue.Set(uFieldValue)
			} else if mFieldValue.Kind() == reflect.Ptr && mFieldValue.Type().Elem() == uFieldValue.Type() {
				ptr := reflect.New(uFieldValue.Type())
				ptr.Elem().Set(uFieldValue)
				mFieldValue.Set(ptr)
			} else if mFieldValue.Kind() == reflect.Ptr && uFieldValue.Kind() == reflect.Ptr && mFieldValue.Type() == uFieldValue.Type() {
				mFieldValue.Set(uFieldValue)
			} else if mFieldValue.Kind() != reflect.Ptr && uFieldValue.Kind() == reflect.Ptr && uFieldValue.Type().Elem() == mFieldValue.Type() {
				if !uFieldValue.IsNil() {
					mFieldValue.Set(uFieldValue.Elem())
				}
			} else {
				return fmt.Errorf("cannot assign %v to %v", uFieldValue.Type(), mFieldValue.Type())
			}
		}
	}

	return nil
}

// updateSlice updates the elements of a slice based on the updates provided.
func updateSlice(slice interface{}, updateSlice interface{}, key string) (interface{}, error) {
	sliceValue := reflect.ValueOf(slice)
	for sliceValue.Kind() == reflect.Ptr || sliceValue.Kind() == reflect.Interface {
		sliceValue = sliceValue.Elem()
	}

	updateSliceValue := reflect.ValueOf(updateSlice)
	for updateSliceValue.Kind() == reflect.Ptr || updateSliceValue.Kind() == reflect.Interface {
		updateSliceValue = updateSliceValue.Elem()
	}

	if sliceValue.Kind() != reflect.Slice || updateSliceValue.Kind() != reflect.Slice {
		return nil, fmt.Errorf("both parameters must be slices")
	}

	elemType := sliceValue.Type().Elem()
	if elemType.Kind() != reflect.Struct {
		return updateSliceValue.Interface(), nil
	}

	resultSlice := reflect.MakeSlice(sliceValue.Type(), 0, sliceValue.Len())

	elemMap := make(map[interface{}]reflect.Value)
	for i := 0; i < sliceValue.Len(); i++ {
		elem := sliceValue.Index(i)
		keyValue, err := getKeyValueFromValue(elem, key)
		if err != nil {
			return nil, err
		}
		elemMap[keyValue] = elem
	}

	for i := 0; i < updateSliceValue.Len(); i++ {
		uElem := updateSliceValue.Index(i)
		keyValue, err := getKeyValueFromValue(uElem, key)
		if err != nil {
			return nil, err
		}

		if elem, exists := elemMap[keyValue]; exists {
			if !elem.CanAddr() {
				return nil, fmt.Errorf("cannot get address of element")
			}
			err := updateModel(elem.Addr().Interface(), uElem.Interface(), key)
			if err != nil {
				return nil, err
			}
			resultSlice = reflect.Append(resultSlice, elem)
			delete(elemMap, keyValue)
		} else {
			newElem, err := convertUpdateToModelInterface(uElem.Interface(), elemType)
			if err != nil {
				return nil, err
			}
			resultSlice = reflect.Append(resultSlice, reflect.ValueOf(newElem))
		}
	}
	return resultSlice.Interface(), nil
}

// getKeyValueFromValue retrieves the key value from a reflect.Value.
func getKeyValueFromValue(v reflect.Value, key string) (interface{}, error) {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	field := v.FieldByName(key)
	if !field.IsValid() {
		return nil, fmt.Errorf("key '%s' not found in struct", key)
	}
	return field.Interface(), nil
}

// isNilOrZero checks if a reflect.Value is nil or zero.
func isNilOrZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return v.IsNil()
	default:
		zero := reflect.Zero(v.Type()).Interface()
		return reflect.DeepEqual(v.Interface(), zero)
	}
}

// convertUpdateToModel converts an update object to a model object.
func convertUpdateToModel[T any](update interface{}) (T, error) {
	var model T
	modelValue := reflect.ValueOf(&model).Elem()
	updateValue := reflect.ValueOf(update)

	for updateValue.Kind() == reflect.Ptr || updateValue.Kind() == reflect.Interface {
		updateValue = updateValue.Elem()
	}

	// Check if updateValue is a struct
	if updateValue.Kind() != reflect.Struct {
		return model, fmt.Errorf("update is not a struct")
	}

	updateType := updateValue.Type()

	for i := 0; i < updateType.NumField(); i++ {
		field := updateType.Field(i)
		uFieldValue := updateValue.Field(i)
		mFieldValue := modelValue.FieldByName(field.Name)
		if !mFieldValue.IsValid() || isNilOrZero(uFieldValue) {
			continue
		}

		if mFieldValue.Kind() == reflect.Struct && uFieldValue.Kind() == reflect.Struct {
			newFieldValue, err := convertUpdateToModelInterface(uFieldValue.Interface(), mFieldValue.Type())
			if err != nil {
				return model, err
			}
			mFieldValue.Set(reflect.ValueOf(newFieldValue))
		} else if mFieldValue.Kind() == reflect.Slice && uFieldValue.Kind() == reflect.Slice {
			uSlice := uFieldValue
			mSliceType := mFieldValue.Type()
			mElemType := mSliceType.Elem()
			mSlice := reflect.MakeSlice(mSliceType, uSlice.Len(), uSlice.Cap())

			for j := 0; j < uSlice.Len(); j++ {
				uElem := uSlice.Index(j)
				newElem, err := convertUpdateToModelInterface(uElem.Interface(), mElemType)
				if err != nil {
					return model, err
				}
				mSlice.Index(j).Set(reflect.ValueOf(newElem))
			}
			mFieldValue.Set(mSlice)
		} else if uFieldValue.Kind() == reflect.Ptr {
			if !uFieldValue.IsNil() {
				mFieldValue.Set(uFieldValue.Elem())
			}
		} else {
			mFieldValue.Set(uFieldValue)
		}
	}

	return model, nil
}

// convertUpdateToModelInterface converts an update object to a model object (interface version).
func convertUpdateToModelInterface(update interface{}, modelType reflect.Type) (interface{}, error) {
	updateValue := reflect.ValueOf(update)
	for updateValue.Kind() == reflect.Ptr {
		updateValue = updateValue.Elem()
	}

	// Check if updateValue is a struct
	if updateValue.Kind() != reflect.Struct {
		// For basic types, check if they are assignable to the modelType
		if updateValue.Type().AssignableTo(modelType) {
			return updateValue.Interface(), nil
		} else {
			return nil, fmt.Errorf("cannot assign %v to %v", updateValue.Type(), modelType)
		}
	}

	modelValue := reflect.New(modelType).Elem()
	updateType := updateValue.Type()

	for i := 0; i < updateType.NumField(); i++ {
		field := updateType.Field(i)
		uFieldValue := updateValue.Field(i)
		mFieldValue := modelValue.FieldByName(field.Name)
		if !mFieldValue.IsValid() || isNilOrZero(uFieldValue) {
			continue
		}

		if mFieldValue.Kind() == reflect.Struct && uFieldValue.Kind() == reflect.Struct {
			newFieldValue, err := convertUpdateToModelInterface(uFieldValue.Interface(), mFieldValue.Type())
			if err != nil {
				return nil, err
			}
			mFieldValue.Set(reflect.ValueOf(newFieldValue))
		} else if mFieldValue.Kind() == reflect.Slice && uFieldValue.Kind() == reflect.Slice {
			uSlice := uFieldValue
			mSliceType := mFieldValue.Type()
			mElemType := mSliceType.Elem()
			mSlice := reflect.MakeSlice(mSliceType, uSlice.Len(), uSlice.Cap())

			for j := 0; j < uSlice.Len(); j++ {
				uElem := uSlice.Index(j)
				newElem, err := convertUpdateToModelInterface(uElem.Interface(), mElemType)
				if err != nil {
					return nil, err
				}
				mSlice.Index(j).Set(reflect.ValueOf(newElem))
			}
			mFieldValue.Set(mSlice)
		} else if uFieldValue.Kind() == reflect.Ptr {
			if !uFieldValue.IsNil() {
				mFieldValue.Set(uFieldValue.Elem())
			}
		} else {
			mFieldValue.Set(uFieldValue)
		}
	}

	return modelValue.Interface(), nil
}
