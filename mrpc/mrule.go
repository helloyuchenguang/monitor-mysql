package mrpc

import (
	cmap "github.com/orcaman/concurrent-map/v2"
	"log/slog"
)

const ruleName = "GRPCRule"

type RuleClient[T any] interface {
	// GetID 获取客户端 ID
	GetID() string
}

// GRPCRule GRPC 规则处理器
type GRPCRule[T any, D any] struct {
	clients cmap.ConcurrentMap[string, *RuleClient[T]]
}

// 自动注册
func init() {
	slog.Info("自动注册 GRPC 规则处理器")

}
