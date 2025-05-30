package main

import (
	"main/common/config"
	"main/mgrpc"
	"main/monitor"
	"main/web"
)

func main() {
	// monitor-mysql
	cfg := config.LoadConfig("./resources/config.yml")
	sseRule, grpcRule := registryRuleService(&cfg)
	monitor.NewMonitorService(cfg, sseRule, grpcRule)
}

func registryRuleService(config *config.Config) (*web.SSERuleService, *mgrpc.GRPCRuleServer) {
	watchHandlers := config.WatchHandlers
	if len(watchHandlers) == 0 {
		return nil, nil
	}
	ruleNameSet := make(map[string]bool, 2)
	var sseRule *web.SSERuleService
	var grpcRule *mgrpc.GRPCRuleServer
	for _, handler := range watchHandlers {
		for _, ruleName := range handler.Rules {
			if _, exists := ruleNameSet[ruleName]; !exists {
				ruleNameSet[ruleName] = true
				switch ruleName {
				case web.RuleName:
					if sseRule == nil {
						sseRule = web.NewSSERuleService(web.Config{
							Addr: config.Web.Addr,
						})
					}
				case mgrpc.RuleName:
					if grpcRule == nil {
						grpcRule = mgrpc.NewGRPCRuleServer(mgrpc.Config{
							Addr: config.GRPC.Addr,
						})
					}
				default:
					panic("未知规则: " + ruleName)
				}
			}
		}
	}
	return sseRule, grpcRule
}
