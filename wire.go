//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"steplems-bot/services"
)
import "steplems-bot/providers"

type WireApplication struct {
	telegramService *services.TelegramService
}

func provideWireApplication(telegramService *services.TelegramService) WireApplication {
	return WireApplication{telegramService: telegramService}
}

func NewWireApplication() (WireApplication, error) {
	wire.Build(provideWireApplication, providers.ProvidersSet, services.ServicesSet)
	return WireApplication{}, nil
}

func (w WireApplication) Start() {
	w.telegramService.StartBot()
}
