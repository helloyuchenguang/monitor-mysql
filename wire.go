//go:build wireinject
// +build wireinject

package monitormysql

import (
	"github.com/google/wire"
	"monitormysql/meili"
)

func InitMeiliService(cfg meili.ClientConfig) *meili.ClientService {
	wire.Build(meili.NewClient, meili.NewMeiliService)
	return nil
}
