//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"main/meili"
)

func InitMeiliService(cfg meili.ClientConfig) *meili.ClientService {
	wire.Build(meili.NewClient, meili.NewMeiliService)
	return nil
}
