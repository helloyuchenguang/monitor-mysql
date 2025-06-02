package main

import (
	"main/common/config"
	"main/mgrpc"
	"main/monitor"
	"main/web"
)

func main() {
	// monitor-mysql
	cfg := config.LoadConfig("./resources/config-13600kf.yml")
	sseRule, grpcRule := registryRuleService(&cfg)
	monitor.InitMonitorService(&cfg, sseRule, grpcRule).StartCanal()
}

func registryRuleService(config *config.Config) (*web.SSERuleService, *mgrpc.GRPCRuleService) {
	watchHandlers := config.WatchHandlers
	if len(watchHandlers) == 0 {
		return nil, nil
	}
	ruleNameSet := make(map[string]bool, 2)
	var sseRule *web.SSERuleService
	var grpcRule *mgrpc.GRPCRuleService
	for _, handler := range watchHandlers {
		for _, ruleName := range handler.Rules {
			if _, exists := ruleNameSet[ruleName]; !exists {
				ruleNameSet[ruleName] = true
				switch ruleName {
				case web.RuleName:
					if sseRule == nil {
						sseRule = web.NewWebSSERuleService(config)
					}
				case mgrpc.RuleName:
					if grpcRule == nil {
						grpcRule = mgrpc.NewGRPCRuleService(config)
					}
				default:
					panic("未知规则: " + ruleName)
				}
			}
		}
	}
	return sseRule, grpcRule
}
