package monitor

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/samber/lo"
	"log/slog"
	"main/common/event"
	"main/common/event/edit"
	"main/common/event/row"
	"main/common/event/rule"
	"regexp"
)

type WatchRegexp struct {
	Regexp *regexp.Regexp
	Rules  []rule.MonitorRuler
}

type CustomEventHandler struct {
	canal.DummyEventHandler
	WatchRegexps []*WatchRegexp
	count        int
}

// isWatched 判断表是否被监控
func (h *CustomEventHandler) isWatched(schema, table string) ([]rule.MonitorRuler, bool) {
	fullName := fmt.Sprintf("%s.%s", schema, table)
	for _, r := range h.WatchRegexps {
		if r.Regexp.MatchString(fullName) {
			return r.Rules, true
		}
	}
	return nil, false
}

// OnRow 处理行事件
func (h *CustomEventHandler) OnRow(e *canal.RowsEvent) error {
	// 判断表是否被监控
	tableSchema := e.Table.Schema
	tableName := e.Table.Name
	rules, ok := h.isWatched(tableSchema, tableName)
	if !ok {
		return nil
	}
	// 根据事件类型生成对应的事件数据
	data := GenerateEventData(e)
	for _, r := range rules {
		// 如果规则没有客户端连接，则跳过
		if r.ClientIsEmpty() {
			continue
		}
		err := r.OnNext(data)
		if err != nil {
			slog.Error(fmt.Sprintf("处理删除事件失败: %v", err))
		}
		h.count++
	}
	return nil
}

// GenerateEventData 根据 canal.RowsEvent 生成对应的 event.Data
func GenerateEventData(e *canal.RowsEvent) *event.Data {
	switch e.Action {
	case canal.InsertAction:
		return &event.Data{
			EventType: event.Insert,
			Table:     event.NewTable(e.Table.Schema, e.Table.Name, e.Table.Columns),
			SaveData: lo.Map(row.ToRowDataList(e.Table.Columns, e.Rows), func(item event.RowData, _ int) *event.SaveData {
				return &event.SaveData{RowData: item}
			}),
		}
	case canal.DeleteAction:
		return &event.Data{
			EventType: event.Delete,
			Table:     event.NewTable(e.Table.Schema, e.Table.Name, e.Table.Columns),
			DeleteData: lo.Map(row.ToRowDataList(e.Table.Columns, e.Rows), func(item event.RowData, _ int) *event.DeleteData {
				return &event.DeleteData{RowData: item}
			}),
		}
	case canal.UpdateAction:
		return &event.Data{
			EventType: event.Update,
			Table:     event.NewTable(e.Table.Schema, e.Table.Name, e.Table.Columns),
			EditData:  edit.ToEditDataList(e.Table.Columns, e.Rows),
		}
	default:
		slog.Warn(fmt.Sprintf("未知的事件类型: %s", e.Action))
		return nil
	}
}
