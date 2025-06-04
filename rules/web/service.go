package web

import (
	"log/slog"
	"main/common/config"
	"main/common/event/rule"
)

const RuleName = "SSERule"

type Config struct {
	Addr string `json:"addr"` // SSE服务地址
}

// SSERuleService rule服务的实体
type SSERuleService struct {
	cfg  *Config
	Rule *rule.RuleServer
}

// NewWebSSERuleService 创建一个新的SSE服务实例
func NewWebSSERuleService(cfg *config.Config) *SSERuleService {
	if !cfg.ExistsRuleName(RuleName) {
		slog.Info("配置中不存在SSERule，不创建SSE服务")
		return nil
	}
	sseRule := SSERuleService{
		cfg:  &Config{Addr: cfg.Web.Addr},
		Rule: rule.NewServer(),
	}
	go sseRule.StartServer()
	slog.Info("注册SSE规则服务", slog.String("addr", sseRule.cfg.Addr))
	return &sseRule
}
