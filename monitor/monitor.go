package monitor

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"log/slog"
	"main/common"
	edit2 "main/common/mevent/edit"
	"regexp"
)

// MyEventHandler 自定义事件处理器
type MyEventHandler struct {
	canal.DummyEventHandler
	WatchRegexps []*regexp.Regexp
	Rules        map[int][]*edit2.MonitorRuler
}

// isWatched 判断表是否被监控
func (h *MyEventHandler) isWatched(schema, table string) ([]*edit2.MonitorRuler, bool) {
	fullName := fmt.Sprintf("%s.%s", schema, table)
	for i, r := range h.WatchRegexps {
		if r.MatchString(fullName) {
			return h.Rules[i], true
		}
	}
	return nil, false
}

// OnRow 处理行事件
func (h *MyEventHandler) OnRow(e *canal.RowsEvent) error {
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
					err := (*rule).OnChange(&edit2.SourceData{
						TableSchema: tableSchema,
						TableName:   tableName,
						Cols:        cols,
						Before:      before,
						After:       after,
					})
					if err != nil {
						slog.Error(fmt.Sprintf("处理更新事件失败: %v", err))
					}
				}()
			}
		}
	}
	return nil
}

// Run 启动
func Run(cfgFile string) (*common.Config, error) {
	// 加载配置文件
	cfg, err := common.LoadConfig(cfgFile)
	if err != nil {
		return nil, err
	}
	// 创建canal.Config
	canalCfg := NewCanalConfigByConfig(&cfg)
	// 创建事件处理器
	handler, err := NewEventHandlerByConfig(&cfg)
	if err != nil {
		slog.Error(fmt.Sprintf("创建事件处理器失败: %v", err))
		return nil, err
	}
	go func() { StartCanal(canalCfg, handler) }()
	return &cfg, nil
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
func NewEventHandlerByConfig(cfg *common.Config) (*MyEventHandler, error) {
	// 把schema和table正则合成一个正则表达式列表给IncludeTableRegex
	var compiledRegexps []*regexp.Regexp
	// 表格正则对应的监控规则
	rules := make(map[int][]*edit2.MonitorRuler, len(cfg.WatchHandlers))
	for i, wt := range cfg.WatchHandlers {
		r, err := regexp.Compile(wt.TableRegex)
		if err != nil {
			slog.Error(fmt.Sprintf("编译正则失败: %v", err))
			return nil, err
		}
		compiledRegexps = append(compiledRegexps, r)
		// 添加规则
		ruleSize := len(wt.Rules)
		// 如果没有规则,使用默认规则
		if ruleSize == 0 {
			slog.Error(fmt.Sprintf("表 %s 没有监控规则,使用默认监控规则", wt.TableRegex))
			rules[i] = []*edit2.MonitorRuler{common.GetDefaultRule()}
			continue
		} else {
			tableRules := make([]*edit2.MonitorRuler, ruleSize)
			for j, ruleName := range wt.Rules {
				rule, ok := common.GetRuleByName(ruleName)
				if !ok {
					slog.Error(fmt.Sprintf("规则 %s 不存在,请检查配置", ruleName))
					return nil, fmt.Errorf("rule %s not found", ruleName)
				}
				tableRules[j] = rule
			}
			rules[i] = tableRules
		}
	}

	eventHandler := &MyEventHandler{
		WatchRegexps: compiledRegexps,
		Rules:        rules,
	}
	return eventHandler, nil
}

// NewCanalConfigByConfig 根据配置文件,创建canal.Config
func NewCanalConfigByConfig(cfg *common.Config) *canal.Config {
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
