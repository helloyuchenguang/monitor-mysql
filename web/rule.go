package web

import (
	"log/slog"
	"main/global"
	"main/global/mevent/edit"
)

const RuleName = "SSERule"

// 自动注册
func init() {
	slog.Info("自动注册 %s 规则处理器", RuleName)
	global.RegisterRule(RuleName, edit.NewServer[[]byte]())
}
