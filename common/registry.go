package common

import (
	edit2 "main/common/mevent/edit"
	"sync"
)

var ruleRegistry = sync.Map{}

// RegisterRule 注册监控处理器
func RegisterRule(name string, handler edit2.MonitorRuler) {
	if _, ok := ruleRegistry.Load(name); ok {
		panic("duplicate register handler: " + name)
	}
	ruleRegistry.Store(name, handler)
}

func GetDefaultRule() *edit2.MonitorRuler {
	if handler, ok := ruleRegistry.Load("SSERule"); ok {
		if rule, ok := handler.(edit2.MonitorRuler); ok {
			return &rule
		}
	}
	return nil
}

// GetRuleByName 根据名称获取监控规则
func GetRuleByName(name string) (*edit2.MonitorRuler, bool) {
	if handler, ok := ruleRegistry.Load(name); ok {
		if rule, ok := handler.(edit2.MonitorRuler); ok {
			return &rule, true
		}
	}
	return nil, false
}

func GetRule[T edit2.ChannelReplyType](ruleName string) *edit2.RuleServer[T] {
	if rule, ok := GetRuleByName(ruleName); ok {
		return (*rule).(*edit2.RuleServer[T])
	}
	return nil
}
