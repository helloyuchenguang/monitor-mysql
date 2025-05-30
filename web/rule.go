package web

import (
	"log/slog"
	"main/common"
	"main/common/mevent/edit"
)

const RuleName = "SSERule"

// 自动注册
func init() {
	slog.Info("自动注册 %s 规则处理器", RuleName)
	common.RegisterRule(RuleName, edit.NewServer[[]byte]())
}
