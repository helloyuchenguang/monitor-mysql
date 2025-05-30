package main

import (
	"main/mgrpc"
	"main/monitor"
	"main/web"
	_ "main/web"
)

func main() {
	// monitor-mysql
	cnf, err := monitor.Run("./properties/properties.yml")
	if err != nil {
		panic(err)
	}
	go mgrpc.RunGrpcCanal(cnf.GRPC.Addr)
	web.StartServer(cnf.Web.Addr)
}
