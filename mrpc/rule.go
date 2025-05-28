package mrpc

import (
	"log/slog"
	"monitormysql/global"
	"monitormysql/global/mevent/edit"
	"monitormysql/mrpc/api/mycanal"
)

const RuleName = "GRPCRule"

// 自动注册
func init() {
	slog.Info("自动注册 %s 规则处理器", RuleName)
	global.RegisterRule(RuleName, edit.NewServer[mycanal.EventTableRowReply]())
}
