package mgrpc

import (
	"log/slog"
	"main/common/config"
	"main/common/event/edit"
)

const RuleName = "GRPCRule"

type Config struct {
	Addr string `json:"addr"` // gRPC服务地址
}

// GRPCRuleService 构建时需要注入的类型
type GRPCRuleService struct {
	cfg  *Config
	Rule *edit.RuleServer
}

// NewGRPCRuleService 创建一个新的GRPCRuleServer实例
func NewGRPCRuleService(cfg *config.Config) *GRPCRuleService {
	if !cfg.ExistsRuleName(RuleName) {
		slog.Info("配置中不存在GRPCRule，不创建GRPC服务")
		return nil
	}
	service := &GRPCRuleService{
		cfg:  &Config{Addr: cfg.GRPC.Addr},
		Rule: edit.NewServer(),
	}
	go service.RunGrpcCanal()
	slog.Info("grpc规则服务已启动", slog.String("addr", service.cfg.Addr))
	return service
}
