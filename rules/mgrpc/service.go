package mgrpc

import (
	"log/slog"
	"main/common/config"
	"main/common/event/rule"
)

const RuleName = "grpc"

type Config struct {
	Enable bool   // 是否启用GRPC规则服务
	Addr   string // gRPC服务地址
}

// GRPCRuleService 构建时需要注入的类型
type GRPCRuleService struct {
	cfg  *Config
	Rule *rule.Server
}

// NewGRPCConfig 创建GRPC服务配置
func NewGRPCConfig(cfg *config.Config) *Config {
	grpcCfg := cfg.SubscribeServerConfig.Grpc
	return &Config{
		Enable: grpcCfg.Enable,
		Addr:   grpcCfg.Addr,
	}
}

// NewGRPCRuleService 创建一个新的GRPCRuleServer实例
func NewGRPCRuleService(cfg *Config) *GRPCRuleService {
	if !cfg.Enable {
		slog.Info("GRPC规则服务未启用", slog.String("addr", cfg.Addr))
		return nil
	}
	service := &GRPCRuleService{
		cfg:  cfg,
		Rule: rule.NewServer(),
	}
	go service.RunGrpcCanal()
	slog.Info("grpc规则服务已启动", slog.String("addr", service.cfg.Addr))
	return service
}
