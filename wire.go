//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"main/common/config"
	"main/monitor"
	"main/rules/meili"
	"main/rules/mgrpc"
	"main/rules/web"
)

// 初始化SSE服务
func InitSSEService(cfg *config.Config) *web.SSERuleService {
	wire.Build(web.NewSSEConfig, web.NewWebSSERuleService)
	return nil
}

// 初始化GRPC服务
func InitGRPCRuleService(cfg *config.Config) *mgrpc.GRPCRuleService {
	wire.Build(mgrpc.NewGRPCConfig, mgrpc.NewGRPCRuleService)
	return nil
}

// 初始化MeiliSearch服务
func InitMeiliService(cfg *config.Config) *meili.ClientService {
	wire.Build(meili.NewMeiliConfig, meili.NewMeiliService)
	return nil
}

// 初始化Monitor服务
func InitMonitorService(cfg *config.Config) *monitor.CanalMonitorService {
	wire.Build(InitSSEService, InitGRPCRuleService, InitMeiliService, monitor.NewMonitorService)
	return nil
}
