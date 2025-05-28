package main

import (
	"monitormysql"
	"monitormysql/mrpc"
	"monitormysql/web"
	_ "monitormysql/web"
)

func main() {
	// monitor-mysql
	cnf, err := monitormysql.Run("./config/config-5900x.yml")
	if err != nil {
		panic(err)
	}
	go mrpc.RunGrpcCanal(cnf.GRPC.Addr)
	web.StartServer(cnf.Web.Addr)
}
