package global

import (
	"monitormysql/global/mevent"
	"sync"
)

var ruleRegistry = sync.Map{}

// RegisterRule 注册监控处理器
func RegisterRule(name string, handler mevent.MonitorRuler) {
	if _, ok := ruleRegistry.Load(name); ok {
		panic("duplicate register handler: " + name)
	}
	ruleRegistry.Store(name, handler)
}

func GetDefaultRule() *mevent.MonitorRuler {
	if handler, ok := ruleRegistry.Load("SSERule"); ok {
		if rule, ok := handler.(mevent.MonitorRuler); ok {
			return &rule
		}
	}
	return nil
}

// GetRuleByName 根据名称获取监控规则
func GetRuleByName(name string) (*mevent.MonitorRuler, bool) {
	if handler, ok := ruleRegistry.Load(name); ok {
		if rule, ok := handler.(mevent.MonitorRuler); ok {
			return &rule, true
		}
	}
	return nil, false
}
