package monitor

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"log/slog"
	"main/common/event/edit"
	"regexp"
)

type CustomEventHandler struct {
	canal.DummyEventHandler
	WatchRegexps []*regexp.Regexp
	Rules        map[int][]edit.MonitorRuler
}

// isWatched 判断表是否被监控
func (h *CustomEventHandler) isWatched(schema, table string) ([]edit.MonitorRuler, bool) {
	fullName := fmt.Sprintf("%s.%s", schema, table)
	for i, r := range h.WatchRegexps {
		if r.MatchString(fullName) {
			return h.Rules[i], true
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

	action := e.Action
	cols := e.Table.Columns
	if len(cols) == 0 {
		slog.Error(fmt.Sprintf("表 %s.%s 没有列信息", tableSchema, tableName))
		return nil
	}

	switch action {
	case canal.UpdateAction:
		for i := 0; i < len(e.Rows); i += 2 {
			before := e.Rows[i]
			after := e.Rows[i+1]
			for _, rule := range rules {
				go func() {
					err := rule.OnChange(edit.ToEditEventData(tableSchema, tableName, cols, before, after))
					if err != nil {
						slog.Error(fmt.Sprintf("处理更新事件失败: %v", err))
					}
				}()
			}
		}
	}
	return nil
}
