package mgrpc

import (
	"log/slog"
	"main/common"
	"main/common/event/edit"
	"main/grpc/api/mycanal"
)

const RuleName = "GRPCRule"

// 自动注册
func init() {
	slog.Info("自动注册 %s 规则处理器", RuleName)
	common.RegisterRule(RuleName, edit.NewServer[*mycanal.EventTableRowReply]())
}
