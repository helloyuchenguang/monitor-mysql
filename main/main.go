package main

import (
	"monitormysql"
	"monitormysql/web"
	_ "monitormysql/web"
)

func main() {
	// monitor-mysql
	monitormysql.Run("./config/config.yml")
	web.StartServer()
}
