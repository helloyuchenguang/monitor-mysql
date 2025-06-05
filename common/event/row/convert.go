package row

import (
	"github.com/go-mysql-org/go-mysql/schema"
	"github.com/samber/lo"
	"main/common/event"
)

// ToRowDataList 转换为rowDataList
func ToRowDataList(cols []schema.TableColumn, rows [][]any) []event.RowData {
	return lo.Map(rows, func(item []any, _ int) event.RowData {
		// 记录行数据
		saveData := event.NewRowDataWithSize(len(cols))
		for idx, col := range cols {
			colVal := item[idx]
			if colVal == nil {
				continue
			}
			// 对于type == []uint8的列，进行转换
			if col.RawType == "text" {
				if val, ok := colVal.([]uint8); ok {
					colVal = string(val)
				}
			}
			saveData[col.Name] = colVal
		}
		return saveData
	})

}
