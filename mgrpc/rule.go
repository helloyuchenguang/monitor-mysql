package mgrpc

import (
	"log/slog"
	"main/common/event/edit"
	"main/mgrpc/api/mycanal"
)

const RuleName = "GRPCRule"

type Config struct {
	Addr string `json:"addr"` // gRPC服务地址
}

// GRPCRuleServer 构建时需要注入的类型
type GRPCRuleServer struct {
	Rule *edit.RuleServer[*mycanal.EventTableRowReply]
}

// NewGRPCRuleServer 创建一个新的GRPCRuleServer实例
func NewGRPCRuleServer(cfg Config) *GRPCRuleServer {
	slog.Info("注册gRPC规则服务", slog.String("addr", cfg.Addr))
	RunGrpcCanal(&cfg)
	grpcRule := &GRPCRuleServer{
		Rule: edit.NewServer[*mycanal.EventTableRowReply](),
	}
	return grpcRule
}
