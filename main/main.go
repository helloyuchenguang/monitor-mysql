package main

import (
	"monitormysql"
	"monitormysql/mgrpc"
	"monitormysql/web"
	_ "monitormysql/web"
)

func main() {
	// monitor-mysql
	cnf, err := monitormysql.Run("./config/config.yml")
	if err != nil {
		panic(err)
	}
	go mgrpc.RunGrpcCanal(cnf.GRPC.Addr)
	web.StartServer(cnf.Web.Addr)
}
