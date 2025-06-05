package event

import (
	"encoding/json"
	"fmt"
	"github.com/go-mysql-org/go-mysql/schema"
	"github.com/samber/lo"
)

// RowData 行数据
type RowData map[string]any

// NewRowDataWithSize 创建一个指定大小的行数据
func NewRowDataWithSize(size int) RowData {
	data := make(RowData, size)
	return data
}

// ConvertMapStrStr 将 RowData 转换为 map[string]string
func (rd RowData) ConvertMapStrStr() map[string]string {
	// 将 RowData 转换为 map[string]string
	result := make(map[string]string, len(rd))
	for k, v := range rd {
		if v == nil {
			continue
		}
		if str, ok := v.(string); ok {
			result[k] = str
		} else {
			result[k] = fmt.Sprintf("%v", v)
		}
	}
	return result
}

// PrimaryKey 获取主键值
func (rd RowData) PrimaryKey() string {
	if id, ok := rd["id"]; ok {
		return fmt.Sprintf("%v", id)
	}
	return ""
}

// EventType 事件类型
type EventType string

const (
	Insert EventType = "INSERT" // 插入事件
	Update EventType = "UPDATE" // 更新事件
	Delete EventType = "DELETE" // 删除事件
)

// Table 表信息
type Table struct {
	Schema  string    `json:"schema"`  // 表所在的数据库
	Table   string    `json:"table"`   // 表名
	Columns []*Column `json:"columns"` // 列信息
}

func (t *Table) ObtainName() string {
	return fmt.Sprintf("%s_%s", t.Schema, t.Table)
}

func (t *Table) ObtainTableName() string {
	return fmt.Sprintf("%s.%s", t.Schema, t.Table)
}

func NewTable(schemaName string, table string, cols []schema.TableColumn) *Table {
	return &Table{
		Schema: schemaName,
		Table:  table,
		Columns: lo.Map(cols, func(item schema.TableColumn, _ int) *Column {
			return &Column{
				Name:    item.Name,
				RowType: item.RawType,
			}
		}),
	}
}

type Column struct {
	Name    string `json:"name"`    // 列名
	RowType string `json:"rowType"` // 列类型
}

type Data struct {
	EventType  EventType     `json:"eventType"`  // 事件类型
	Table      *Table        `json:"table"`      // 事件对应的表信息
	SaveData   []*SaveData   `json:"saveData"`   // 插入事件数据
	DeleteData []*DeleteData `json:"deleteData"` // 删除事件数据
	EditData   []*EditData   `json:"editData"`   // 编辑事件数据
}

func (data *Data) ToJson() string {
	// 将 Data 转换为 JSON 字符串
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "{}"
	}
	return string(jsonData)
}

// EditData 编辑事件数据
type EditData struct {
	UnChangeRowData RowData                    `json:"unChangeRowData"` // 修改前的行数据
	EditFieldValues map[string]*EditFieldValue `json:"editFieldValues"`
}

func NewEditDataWithSize(size int) *EditData {
	return &EditData{
		UnChangeRowData: NewRowDataWithSize(size),
		EditFieldValues: map[string]*EditFieldValue{},
	}
}

type EditFieldValue struct {
	Before any `json:"before"` // 修改前的值
	After  any `json:"after"`  // 修改后的值
}

func (ef *EditFieldValue) ConvertToPairStr() (string, string) {
	return fmt.Sprintf("%v", ef.Before), fmt.Sprintf("%v", ef.After)
}

// SaveData 插入事件数据
type SaveData struct {
	RowData RowData `json:"rowData"` // 插入的行数据
}

// DeleteData 删除事件数据
type DeleteData struct {
	RowData RowData `json:"rowData"` // 插入的行数据
}
