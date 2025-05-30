package edit

// MonitorRuler 监控处理器接口
type MonitorRuler interface {
	// OnChange 处理更新信息
	OnChange(t *SourceData) error
}
