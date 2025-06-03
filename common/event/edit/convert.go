package edit

import (
	"github.com/go-mysql-org/go-mysql/schema"
	"main/common/event"
	"reflect"
)

// ToEditEventData 从行数据中获取更新信息
func ToEditEventData(tableSchema, tableName string, cols []schema.TableColumn, before, after []any) *event.Data {
	// 更新列
	// 列信息
	var columns []*event.Column
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

		// 封装列信息
		column := event.Column{
			Name:    colName,
			RowType: col.RawType,
		}
		columns = append(columns, &column)

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
	return &event.Data{
		EventType: event.Update,
		Table:     event.NewTable(tableSchema, tableName, columns),
		EditData:  editData,
	}

}
