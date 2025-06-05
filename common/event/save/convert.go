package save

import (
	"github.com/go-mysql-org/go-mysql/schema"
	"main/common/event"
)

func ToSaveEventData(schema, tableName string, cols []schema.TableColumn, rows [][]any) *event.Data {
	// 封装数据
	saveDataList := make([]*event.SaveData, len(rows))
	for i, row := range rows {
		// 记录行数据
		saveData := event.NewRowDataWithSize(len(cols))
		for idx, col := range cols {
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
			saveData[col.Name] = colVal
		}
		saveDataList[i] = &event.SaveData{
			RowData: saveData,
		}
	}
	return &event.Data{
		EventType: event.Insert,
		Table:     event.NewTable(schema, tableName, cols),
		SaveData:  saveDataList,
	}

}
