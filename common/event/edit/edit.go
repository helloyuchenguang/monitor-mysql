package edit

import "main/common/event"

// MonitorRuler 监控处理器接口
type MonitorRuler interface {
	// OnNext 处理更新信息
	OnNext(d *event.Data) error
}
