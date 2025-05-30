package main

import (
	"main/grpc"
	"main/monitor"
	"main/web"
	_ "main/web"
)

func main() {
	// monitor-mysql
	cnf, err := monitor.Run("./resources/resources.yml")
	if err != nil {
		panic(err)
	}
	go grpc.RunGrpcCanal(cnf.GRPC.Addr)
	web.StartServer(cnf.Web.Addr)
}
