package main

import (
	"main/common/config"
)

func main() {
	cfg := config.LoadConfig("./resources/config.yml")
	// monitor-mysql
	InitMonitorService(cfg).StartCanal()
}
