package monitormysql

var TableRegistry map[string]MonitorHandler

func init() {
	TableRegistry = make(map[string]MonitorHandler)
	TableRegistry["TplNodeHandler"] = &TplNodeHandler{
		SseServer: NewSSEServer(),
	}
}
