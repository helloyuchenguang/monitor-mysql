package web

import (
	"log/slog"
	"main/common/config"
	"main/common/event/edit"
)

const RuleName = "SSERule"

type Config struct {
	Addr string `json:"addr"` // SSE服务地址
}

// SSERuleService rule服务的实体
type SSERuleService struct {
	cfg  *Config
	Rule *edit.RuleServer[[]byte]
}

// NewWebSSERuleService 创建一个新的SSE服务实例
func NewWebSSERuleService(cfg *config.Config) *SSERuleService {
	sseRule := SSERuleService{
		cfg:  &Config{Addr: cfg.Web.Addr},
		Rule: edit.NewServer[[]byte](),
	}
	go sseRule.StartServer()
	slog.Info("注册SSE规则服务", slog.String("addr", sseRule.cfg.Addr))
	return &sseRule
}
