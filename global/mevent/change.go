package mevent

import (
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/schema"
	"reflect"
	"regexp"
)

// MonitorRule 监控处理器接口
type MonitorRule interface {
	// OnChange 处理更新信息
	OnChange(ui UpdateInfo) error
}

// MyEventHandler 自定义事件处理器
type MyEventHandler struct {
	canal.DummyEventHandler
	WatchRegexps []*regexp.Regexp
	Handlers     []MonitorRule
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

// FromRows 从行数据中获取更新信息
func FromRows(tableSchema, tableName string,
	cols []schema.TableColumn,
	before, after []any) UpdateInfo {
	// 更新列
	edits := make(map[string]UpdateValueInfo)
	// 列信息
	columnMap := make(map[string]ColumnInfo, len(cols))
	// 记录旧值和新值
	dataUpdate := DataUpdateInfo{
		Old: make(map[string]any),
		New: make(map[string]any),
	}
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
	ui := UpdateInfo{
		TableSchema: tableSchema,
		TableName:   tableName,
		DataUpdate:  dataUpdate,
		ColumnMap:   columnMap,
		Edits:       edits,
	}
	return ui
}
