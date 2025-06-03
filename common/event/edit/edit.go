package edit

import "main/common/event"

// MonitorRuler 监控处理器接口
type MonitorRuler interface {
	// OnChange 处理更新信息
	OnChange(d *event.Data) error
}
