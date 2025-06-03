package del

import (
	"github.com/go-mysql-org/go-mysql/schema"
	"main/common/event"
)

func ToDeleteEventData(schema, tableName string, cols []schema.TableColumn, row []any) *event.Data {
	// 列信息
	var columns []*event.Column
	// 记录行数据
	deleteData := event.NewRowDataWithSize(len(cols))
	for idx, col := range cols {
		colName := col.Name
		colVal := row[idx]
		if colVal == nil {
			continue
		}
		// 对于type == []uint8的列，进行转换
		if col.RawType == "text" {
			if val, ok := colVal.([]uint8); ok {
				colVal = string(val)
			}
		}
		// 封装列信息
		column := event.Column{
			Name:    colName,
			RowType: col.RawType,
		}
		columns = append(columns, &column)
		deleteData[colName] = colVal
	}

	return &event.Data{
		EventType: event.Delete,
		Table:     event.NewTable(schema, tableName, columns),
		DeleteData: &event.DeleteData{
			RowData: deleteData,
		},
	}
}
