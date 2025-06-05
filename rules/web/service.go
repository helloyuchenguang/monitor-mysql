package web

import (
	"log/slog"
	"main/common/config"
	"main/common/event/rule"
)

const RuleName = "sse"

type Config struct {
	Enable bool
	Addr   string // SSE服务地址
}

// SSERuleService rule服务的实体
type SSERuleService struct {
	cfg  *Config
	Rule *rule.Server
}

func NewSSEConfig(cfg *config.Config) *Config {
	sseCfg := cfg.SubscribeServerConfig.SSE
	return &Config{
		Enable: sseCfg.Enable,
		Addr:   sseCfg.Addr,
	}
}

// NewWebSSERuleService 创建一个新的SSE服务实例
func NewWebSSERuleService(cfg *Config) *SSERuleService {
	if !cfg.Enable {
		slog.Info("SSE规则服务未启用", slog.String("addr", cfg.Addr))
		return nil
	}
	sseRule := SSERuleService{
		cfg:  cfg,
		Rule: rule.NewServer(),
	}
	go sseRule.StartServer()
	slog.Info("注册SSE规则服务", slog.String("addr", sseRule.cfg.Addr))
	return &sseRule
}
