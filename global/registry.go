package global

import (
	"monitormysql/global/mevent"
	"sync"
)

var ruleRegistry = sync.Map{}

// RegisterRule 注册监控处理器
func RegisterRule(name string, handler mevent.MonitorRule) {
	if _, ok := ruleRegistry.Load(name); ok {
		panic("duplicate register handler: " + name)
	}
	ruleRegistry.Store(name, handler)
}

// GetAllRules 获取指定名称的监控处理器
func GetAllRules() map[string]mevent.MonitorRule {
	rules := make(map[string]mevent.MonitorRule)
	ruleRegistry.Range(func(key, value any) bool {
		if handler, ok := value.(mevent.MonitorRule); ok {
			rules[key.(string)] = handler
		}
		return true
	})
	return rules
}

func GetDefaultRule() *mevent.MonitorRule {
	if handler, ok := ruleRegistry.Load("SSERule"); ok {
		if rule, ok := handler.(mevent.MonitorRule); ok {
			return &rule
		}
	}
	return nil
}

// GetRuleByName 根据名称获取监控规则
func GetRuleByName(name string) (*mevent.MonitorRule, bool) {
	if handler, ok := ruleRegistry.Load(name); ok {
		if rule, ok := handler.(mevent.MonitorRule); ok {
			return &rule, true
		}
	}
	return nil, false
}
