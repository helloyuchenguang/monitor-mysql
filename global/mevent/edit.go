package mevent

import (
	"github.com/go-mysql-org/go-mysql/canal"
	"monitormysql/mrpc/api/mycanal"
	"regexp"
)

type EditReplyType interface {
	[]byte | mycanal.EventTableRowReply
}

// MonitorRuler 监控处理器接口
type MonitorRuler interface {
	// OnChange 处理更新信息
	OnChange(t *EditSourceData) error
}

// MyEventHandler 自定义事件处理器
type MyEventHandler struct {
	canal.DummyEventHandler
	WatchRegexps []*regexp.Regexp
	Handlers     []MonitorRuler
}
