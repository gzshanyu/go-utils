package json

import (
	"fmt"
	"reflect"
)

func DeepFields(ifaceType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	for i := 0; i < ifaceType.NumField(); i++ {
		v := ifaceType.Field(i)
		if v.Anonymous && v.Type.Kind() == reflect.Struct {
			fields = append(fields, DeepFields(v.Type)...)
		} else {
			fields = append(fields, v)
		}
	}

	return fields
}

func StructCopy(DstStructPtr interface{}, SrcStructPtr interface{}) {
	srcv := reflect.ValueOf(SrcStructPtr)
	dstv := reflect.ValueOf(DstStructPtr)
	srct := reflect.TypeOf(SrcStructPtr)
	dstt := reflect.TypeOf(DstStructPtr)

	if srct.Kind() != reflect.Ptr || dstt.Kind() != reflect.Ptr ||
		srct.Elem().Kind() == reflect.Ptr || dstt.Elem().Kind() == reflect.Ptr {
		panic("Fatal error:type of parameters must be Ptr of value")
	}
	if srcv.IsNil() || dstv.IsNil() {
		panic("Fatal error:value of parameters should not be nil")
	}
	srcV := srcv.Elem()
	dstV := dstv.Elem()

	srcfields := DeepFields(reflect.ValueOf(SrcStructPtr).Elem().Type())
	for _, v := range srcfields {
		if v.Anonymous {
			continue
		}
		dst := dstV.FieldByName(v.Name)
		src := srcV.FieldByName(v.Name)
		if !dst.IsValid() {
			continue
		}

		if src.Type() == dst.Type() && dst.CanSet() {
			dst.Set(src)
			continue
		}

		if src.Kind() == reflect.Ptr && !src.IsNil() && src.Type().Elem() == dst.Type() {
			dst.Set(src.Elem())
			continue
		}
		if dst.Kind() == reflect.Ptr && dst.Type().Elem() == src.Type() {
			dst.Set(reflect.New(src.Type()))
			dst.Elem().Set(src)
			continue
		}

	}
	return
}

// 结构体转为Map[string]interface{}
func Struct2Map(in interface{}, tagName string) (interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // 非结构体返回错误提示
		return nil, fmt.Errorf("ToMap only accepts gStruct or gStruct pointer; got %T", v)
	}

	t := v.Type()
	// 遍历结构体字段
	// 指定tagName值为map中key;字段值为map中value
	if tagName == "" {
		tagName = "json"
	}
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get(tagName); tagValue != "" {
			out[tagValue] = v.Field(i).Interface()
		}
	}
	return out, nil
}

// 获取结构体中某标签的值集合
func GetTagValue(in interface{}, tagName string) ([]interface{}, error) {
	out := make([]interface{}, 0)

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // 非结构体返回错误提示
		return nil, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	t := v.Type()
	// 遍历结构体字段
	// 指定tagName值为map中key;字段值为map中value
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get(tagName); tagValue != "" {
			out = append(out, tagValue)
		}
	}
	return out, nil
}

// 获取结构体字段的值集合
func GetValue(in interface{}, tagName string) ([]interface{}, error) {
	out := make([]interface{}, 0)

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // 非结构体返回错误提示
		return nil, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	t := v.Type()
	// 遍历结构体字段
	// 指定tagName值为map中key;字段值为map中value
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get(tagName); tagValue != "" {
			out = append(out, v.Field(i).Interface())
		}
	}
	return out, nil
}

// 数组是否包含某元素（string）
func IsContainString(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// 数组是否包含某元素（int32）
func IsContainInt32(items []int32, item int32) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// 数组是否包含某元素（int64）
func IsContainInt64(items []int64, item int64) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// 二分法反转数组
func Reverse(arr *[]string, length int) {
	var temp string
	for i := 0; i < length/2; i++ {
		temp = (*arr)[i]
		(*arr)[i] = (*arr)[length-1-i]
		(*arr)[length-1-i] = temp
	}
}

// 数字转换为Excel列A-Z
func ConvertNumToExcelChars(num int) string {
	var (
		excelChar []string
		cols      string
		v         int
	)

	excelChar = []string{"", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	v = num
	for v > 0 {
		k := v % 26
		if k == 0 {
			k = 26
		}
		v = (v - k) / 26
		cols = excelChar[k] + cols
	}

	return cols
}

