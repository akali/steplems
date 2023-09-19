//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"steplems-bot/persistence"
	"steplems-bot/services"
	"steplems-bot/services/spotify"
	"steplems-bot/services/telegram"
)
import "steplems-bot/providers"

type WireApplication struct {
	telegramService *telegram.TelegramService
	ss              *spotify.SpotifyService
}

func provideWireApplication(telegramService *telegram.TelegramService, service *spotify.SpotifyService) WireApplication {
	return WireApplication{telegramService: telegramService, ss: service}
}

func NewWireApplication() (WireApplication, error) {
	wire.Build(provideWireApplication,
		providers.ProvidersSet,
		services.ServicesSet,
		persistence.PersistenceSet)
	return WireApplication{}, nil
}

func (w WireApplication) Start() {
	w.telegramService.StartBot()
}
