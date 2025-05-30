package web

import (
	"log/slog"
	"main/common/event/edit"
)

const RuleName = "SSERule"

type Config struct {
	Addr string `json:"addr"` // SSE服务地址
}

// SSERuleService rule服务的实体
type SSERuleService struct {
	Rule *edit.RuleServer[[]byte]
}

// NewSSERuleService 创建一个新的SSE服务实例
func NewSSERuleService(cfg Config) *SSERuleService {
	slog.Info("注册SSE规则服务", slog.String("addr", cfg.Addr))
	StartServer(&cfg)
	return &SSERuleService{Rule: edit.NewServer[[]byte]()}
}
