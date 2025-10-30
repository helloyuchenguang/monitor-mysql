package edit

import (
	"main/common/event"
	"reflect"

	"github.com/go-mysql-org/go-mysql/schema"
)

// ToEditDataList 从行数据中获取更新信息
func ToEditDataList(cols []schema.TableColumn, rows [][]any) []*event.EditData {
	var editDataList []*event.EditData
	for i := 0; i < len(rows); i += 2 {
		before := rows[i]  // 旧数据
		after := rows[i+1] // 新数据
		// 记录旧值和新值
		editData := event.NewEditDataWithSize(len(cols))
		for idx, col := range cols {
			oldVal := before[idx]
			newVal := after[idx]
			colName := col.Name

			if oldVal == nil && newVal == nil {
				continue
			}
			// 对于type == []uint8的列，进行转换
			colValType := reflect.TypeOf(oldVal)
			switch colValType {
			case reflect.TypeOf([]uint8{}):
				if col.RawType == "text" {
					oldVal = string(oldVal.([]uint8))
					newVal = string(newVal.([]uint8))
				}
			}
			// 如果旧值和新值不相等，则记录更新信息
			if oldVal != newVal {
				editData.EditFieldValues[colName] = &event.EditFieldValue{
					Before: oldVal,
					After:  newVal,
				}
			} else {
				// 记录未更改的值
				if newVal != nil {
					editData.UnChangeRowData[col.Name] = newVal
				}
			}
		}
		editDataList = append(editDataList, editData)
	}
	return editDataList

}
