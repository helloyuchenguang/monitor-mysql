package monitor

import (
	"fmt"
	"main/common/mevent/edit"
	"reflect"
	"strings"
)

func MapToStructByReflect(data map[string]any, out any) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("out must be a pointer to struct")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := t.Field(i)

		// 忽略无法设置的字段（如未导出字段）
		if !field.CanSet() {
			continue
		}

		// 获取字段名和 json tag
		name := structField.Name
		tag := structField.Tag.Get("column")
		if tag != "" && tag != "-" {
			// 只取 json tag 的主名，去掉 ",omitempty" 等
			name = strings.Split(tag, ",")[0]
		}

		if val, ok := data[name]; ok {
			valValue := reflect.ValueOf(val)

			// 尝试类型转换或赋值
			if valValue.Type().AssignableTo(field.Type()) {
				field.Set(valValue)
			} else if valValue.Type().ConvertibleTo(field.Type()) {
				field.Set(valValue.Convert(field.Type()))
			} else {
				// 跳过无法赋值的字段（类型不匹配）
				continue
			}
		}
	}

	return nil
}

// ConvertByUpdateInfo 将 UpdateInfo 转换为指定类型的结构体
func ConvertByUpdateInfo[T any](updateInfo edit.UpdateInfo) (T, T, error) {
	var oldRaw T
	var newRaw T
	err := MapToStructByReflect(updateInfo.DataUpdate.Old, &oldRaw)
	if err != nil {
		return oldRaw, newRaw, err
	}
	err = MapToStructByReflect(updateInfo.DataUpdate.New, &newRaw)
	if err != nil {
		return oldRaw, newRaw, err
	}
	return oldRaw, newRaw, nil
}
