//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"main/meili"
	"main/mgrpc"
	"main/web"
)

// 初始化SSE服务
func InitSSEService(cfg web.Config) *web.SSERuleService {
	wire.Build(web.NewSSERuleService)
	return nil
}

// 初始化SSE服务
func InitGRPCRuleService(cfg mgrpc.Config) *mgrpc.GRPCRuleServer {
	wire.Build(mgrpc.NewGRPCRuleServer)
	return nil
}

// 初始化MeiliSearch服务
func InitMeiliService(cfg meili.ClientConfig) *meili.ClientService {
	wire.Build(meili.NewClient, meili.NewMeiliService)
	return nil
}
