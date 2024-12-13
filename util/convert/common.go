package convert

import (
	"reflect"
	"time"
)

func CreateFieldConvert(srcVal, dstVal reflect.Value) {
	// 若型別可直接指派
	if dstVal.IsValid() && srcVal.Type().AssignableTo(dstVal.Type()) && dstVal.CanSet() {
		dstVal.Set(srcVal)
		return
	}

	switch srcVal.Kind() {
	case reflect.Ptr:
		// 處理指標
		if srcVal.IsNil() {
			return
		}
		if dstVal.Kind() == reflect.Ptr {
			if dstVal.IsNil() {
				newDst := reflect.New(dstVal.Type().Elem())
				dstVal.Set(newDst)
			}
			CreateFieldConvert(srcVal.Elem(), dstVal.Elem())
		} else {
			// 來源是指標但目標不是指標，則解引用後再處理
			CreateFieldConvert(srcVal.Elem(), dstVal)
		}
	case reflect.Struct:
		// 處理 struct 中的每個匯出欄位
		for i := 0; i < srcVal.NumField(); i++ {
			field := srcVal.Type().Field(i)
			if field.PkgPath != "" {
				// 跳過未導出的欄位
				continue
			}
			srcField := srcVal.Field(i)
			dstField := dstVal.FieldByName(field.Name)
			if dstField.IsValid() && dstField.CanSet() {
				CreateFieldConvert(srcField, dstField)
			}
		}
	case reflect.Slice, reflect.Array:
		// 處理 slice/array
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
		// 為基礎型別嘗試指派或做型別轉換
		if !dstVal.CanSet() {
			return
		}

		srcKind := srcVal.Kind()
		dstKind := dstVal.Kind()

		// 特別處理從 string -> *string
		if srcKind == reflect.String && dstKind == reflect.Ptr && dstVal.Type().Elem().Kind() == reflect.String {
			// 建立 *string
			str := srcVal.String()
			dstVal.Set(reflect.ValueOf(&str))
			return
		}

		// 特別處理從 string -> *time.Time (需要自行定義日期格式)
		if srcKind == reflect.String && dstKind == reflect.Ptr && dstVal.Type().Elem() == reflect.TypeOf(time.Time{}) {
			str := srcVal.String()
			// 假設日期格式為 "2006-01-02"
			parsedTime, err := time.ParseInLocation("2006-01-02", str, time.UTC)
			if err == nil {
				dstVal.Set(reflect.ValueOf(&parsedTime))
			}
			return
		}

		// 特別處理從 string -> *xxx (若有其他型別需求自行擴充)
		// ... (可依需求擴充)

		// 尝试直接设置
		if srcVal.Type().AssignableTo(dstVal.Type()) {
			dstVal.Set(srcVal)
		}
	}
}
