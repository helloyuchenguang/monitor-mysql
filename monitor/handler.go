package monitor

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"log/slog"
	"main/common/event"
	"main/common/event/del"
	"main/common/event/edit"
	"main/common/event/rule"
	"main/common/event/save"
	"regexp"
)

type WatchRegexp struct {
	Regexp *regexp.Regexp
	Rules  []rule.MonitorRuler
}

type CustomEventHandler struct {
	canal.DummyEventHandler
	WatchRegexps []*WatchRegexp
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

	cols := e.Table.Columns
	if len(cols) == 0 {
		slog.Error(fmt.Sprintf("表 %s.%s 没有列信息", tableSchema, tableName))
		return nil
	}
	// 根据事件类型生成对应的事件数据
	dataList := GenerateEventDataList(e)
	for _, r := range rules {
		// 如果规则没有客户端连接，则跳过
		if r.ClientIsEmpty() {
			continue
		}
		err := r.OnNext(dataList)
		if err != nil {
			slog.Error(fmt.Sprintf("处理删除事件失败: %v", err))
		}
	}
	return nil
}

func GenerateEventDataList(e *canal.RowsEvent) *event.Data {
	switch e.Action {
	case canal.InsertAction:
		return save.ToSaveEventData(e.Table.Schema, e.Table.Name, e.Table.Columns, e.Rows)
	case canal.DeleteAction:
		return del.ToDeleteEventData(e.Table.Schema, e.Table.Name, e.Table.Columns, e.Rows)
	case canal.UpdateAction:
		return edit.ToEditEventData(e.Table.Schema, e.Table.Name, e.Table.Columns, e.Rows)
	default:
		slog.Warn(fmt.Sprintf("未知的事件类型: %s", e.Action))
		return nil
	}
}
