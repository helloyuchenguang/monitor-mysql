package edit

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/schema"
	"monitormysql/mgrpc/api/mycanal"
	"reflect"
)

type SourceData struct {
	TableSchema, TableName string
	Cols                   []schema.TableColumn
	Before, After          []any
}

// SourceDataToEditInfo 从行数据中获取更新信息
func SourceDataToEditInfo(sd *SourceData) *UpdateInfo {
	// 更新列
	edits := make(map[string]UpdateValueInfo)
	cols := sd.Cols
	// 列信息
	columnMap := make(map[string]ColumnInfo, len(cols))
	// 记录旧值和新值
	dataUpdate := DataUpdateInfo{
		Old: make(map[string]any),
		New: make(map[string]any),
	}
	before := sd.Before
	after := sd.After
	for idx, col := range cols {
		oldVal := before[idx]
		newVal := after[idx]
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
		colName := col.Name
		// 封装列信息
		columnMap[colName] = ColumnInfo{
			Name:    col.Name,
			RowType: col.RawType,
		}
		// 记录旧值和新值
		if oldVal != nil {
			dataUpdate.Old[col.Name] = oldVal
		}
		if newVal != nil {
			dataUpdate.New[col.Name] = newVal
		}
		// 如果旧值和新值不相等，则记录更新信息
		if oldVal != newVal {
			edits[colName] = UpdateValueInfo{
				Old: oldVal,
				New: newVal,
			}
		}
	}
	return &UpdateInfo{
		TableSchema: sd.TableSchema,
		TableName:   sd.TableName,
		DataUpdate:  dataUpdate,
		ColumnMap:   columnMap,
		Edits:       edits,
	}

}

// SourceDataToGrpcReply 从行数据中获取更新信息
func SourceDataToGrpcReply(sd *SourceData) *mycanal.EventTableRowReply {
	cols := sd.Cols
	// 列信息
	columns := make([]*mycanal.ColumnInfoReply, len(cols))
	// 记录旧值和新值
	editColumns := make(map[string]*mycanal.ColumnValueReply)
	dataUpdate := mycanal.EventUpdateInfoReply{
		BeforeRowData: make(map[string]string),
		AfterRowData:  make(map[string]string),
		EditColumns:   editColumns,
	}
	before := sd.Before
	after := sd.After
	for idx, col := range cols {
		oldVal := before[idx]
		newVal := after[idx]
		if oldVal == nil && newVal == nil {
			continue
		}
		// 封装列信息
		columns = append(columns, &mycanal.ColumnInfoReply{
			ColumnName: col.Name,
			ColumnType: col.RawType,
		})
		// 对于type == []uint8的列，进行转换
		colValType := reflect.TypeOf(oldVal)
		switch colValType {
		case reflect.TypeOf([]uint8{}):
			if col.RawType == "text" {
				oldVal = string(oldVal.([]uint8))
				newVal = string(newVal.([]uint8))
			}
		}
		colName := col.Name
		// 封装列信息
		// 记录旧值和新值
		if oldVal != nil {
			dataUpdate.BeforeRowData[col.Name] = fmt.Sprintf("%v", oldVal)
		}
		if newVal != nil {
			dataUpdate.AfterRowData[col.Name] = fmt.Sprintf("%v", newVal)
		}
		// 如果旧值和新值不相等，则记录更新信息
		if oldVal != newVal {
			editColumns[colName] = &mycanal.ColumnValueReply{
				BeforeValue: fmt.Sprintf("%v", oldVal),
				AfterValue:  fmt.Sprintf("%v", newVal),
			}
		}
	}
	return &mycanal.EventTableRowReply{
		EventType: mycanal.CanalEventType_UPDATE,
		RowStruct: &mycanal.TableStructReply{
			Schema: sd.TableSchema,
			Table:  sd.TableName,
			//Columns: columns,
		},
		UpdateInfo: &dataUpdate,
	}
}

// UpdateInfo 更新信息
type UpdateInfo struct {
	TableSchema string
	TableName   string
	ColumnMap   map[string]ColumnInfo `json:"-"`
	// 数据更新信息
	DataUpdate DataUpdateInfo
	// 更新的列信息
	Edits map[string]UpdateValueInfo
}

type DataUpdateInfo struct {
	Old, New map[string]any
}

// ColumnInfo 列信息
type ColumnInfo struct {
	Name, RowType string
}

// UpdateValueInfo 更新值信息
type UpdateValueInfo struct {
	Old, New any
}
