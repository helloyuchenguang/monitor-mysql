package monitormysql

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"log/slog"
	"regexp"
)

// MyEventHandler 自定义事件处理器
type MyEventHandler struct {
	canal.DummyEventHandler
	WatchRegexps []*regexp.Regexp
	Handlers     []MonitorHandler
}

// isWatched 判断表是否被监控
func (h *MyEventHandler) isWatched(schema, table string) (MonitorHandler, bool) {
	if len(h.WatchRegexps) == 0 {
		return nil, true
	}
	fullName := fmt.Sprintf("%s.%s", schema, table)
	for i, r := range h.WatchRegexps {
		if r.MatchString(fullName) {
			handler := h.Handlers[i]
			return handler, true
		}
	}
	return nil, false
}

// OnRow 处理行事件
func (h *MyEventHandler) OnRow(e *canal.RowsEvent) error {
	// 判断表是否被监控
	tableSchema := e.Table.Schema
	tableName := e.Table.Name
	handler, ok := h.isWatched(tableSchema, tableName)
	if !ok {
		return nil
	}

	action := e.Action
	cols := e.Table.Columns

	switch action {
	case canal.UpdateAction:
		for i := 0; i < len(e.Rows); i += 2 {
			before := e.Rows[i]
			after := e.Rows[i+1]
			// 转换为UpdateInfo
			ui := FromRows(tableSchema, tableName, cols, before, after)
			// 处理更新信息
			err := handler.OnChange(ui)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Run 启动
func Run(cfgFile string) {
	// 加载配置文件
	cfg, err := LoadConfig(cfgFile)
	if err != nil {
		return
	}
	// 创建canal.Config
	canalCfg := NewCanalConfigByConfig(&cfg)
	// 创建事件处理器
	handler, err := NewEventHandlerByConfig(&cfg)
	if err != nil {
		slog.Error(fmt.Sprintf("创建事件处理器失败: %v", err))
		return
	}
	go func() { StartCanal(canalCfg, handler) }()
	return
}

// StartCanal 启动canal
func StartCanal(canalCfg *canal.Config, handler *MyEventHandler) {
	// 创建canal实例
	c, err := canal.NewCanal(canalCfg)
	if err != nil {
		slog.Error(fmt.Sprintf("创建canal失败: %v", err))
		return
	}
	// 设置事件处理器
	c.SetEventHandler(handler)

	// 获取当前的binlog位置
	pos, err := c.GetMasterPos()
	if err != nil {
		slog.Error(fmt.Sprintf("获取masterPos失败: %v", err))
		return
	}

	if err := c.RunFrom(pos); err != nil {
		slog.Error(fmt.Sprintf("canal运行失败: %v", err))
	}
}

// NewEventHandlerByConfig 根据配置文件,创建事件处理器
func NewEventHandlerByConfig(cfg *Config) (*MyEventHandler, error) {
	// 把schema和table正则合成一个正则表达式列表给IncludeTableRegex
	var compiledRegexps []*regexp.Regexp
	var handlers []MonitorHandler
	for _, wt := range cfg.WatchHandlers {
		r, err := regexp.Compile(wt.TableRegex)
		if err != nil {
			slog.Error(fmt.Sprintf("编译正则失败: %v", err))
			return nil, err
		}
		compiledRegexps = append(compiledRegexps, r)
		// 校验handler是否存在
		if handler, ok := TableRegistry[wt.Handler]; ok {
			handlers = append(handlers, handler)
		} else {
			slog.Error(fmt.Sprintf("handler不存在: %v", wt.Handler))
			panic(fmt.Errorf("handler不存在: %v", wt.Handler))
		}
	}
	eventHandler := &MyEventHandler{
		WatchRegexps: compiledRegexps,
		Handlers:     handlers,
	}
	return eventHandler, nil
}

// NewCanalConfigByConfig 根据配置文件,创建canal.Config
func NewCanalConfigByConfig(cfg *Config) *canal.Config {
	canalCfg := canal.NewDefaultConfig()
	canalCfg.Addr = cfg.Database.Addr
	canalCfg.User = cfg.Database.User
	canalCfg.Password = cfg.Database.Password
	canalCfg.Flavor = cfg.Database.Flavor
	canalCfg.ServerID = cfg.Database.ServerID
	canalCfg.Dump.ExecutionPath = cfg.Database.DumpExecutionPath
	canalCfg.IncludeTableRegex = cfg.Database.IncludeTableRegex
	return canalCfg
}
