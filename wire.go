//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"main/common/config"
	"main/meili"
	"main/mgrpc"
	"main/web"
)

// 初始化SSE服务
func InitSSEService(cfg *config.Config) *web.SSERuleService {
	wire.Build(web.NewWebSSERuleService)
	return nil
}

// 初始化GRPC服务
func InitGRPCRuleService(cfg *config.Config) *mgrpc.GRPCRuleService {
	wire.Build(mgrpc.NewGRPCRuleService)
	return nil
}

// 初始化MeiliSearch服务
func InitMeiliService(cfg meili.ClientConfig) *meili.ClientService {
	wire.Build(meili.NewClient, meili.NewMeiliService)
	return nil
}

// 初始化Monitor服务

func InitMonitorService(cfg *config.Config) (*web.SSERuleService, *mgrpc.GRPCRuleService) {
	wire.Build(web.NewWebSSERuleService, mgrpc.NewGRPCRuleService, registryRuleService)
	return nil, nil
}
