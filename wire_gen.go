// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"steplems-bot/providers"
	"steplems-bot/services"
)

// Injectors from wire.go:

func NewWireApplication() (WireApplication, error) {
	telegramBotToken, err := providers.ProvideBotToken()
	if err != nil {
		return WireApplication{}, err
	}
	telegramWebhookAddress, err := providers.ProvideBotWebhook()
	if err != nil {
		return WireApplication{}, err
	}
	botAPI, err := providers.ProvideTgBot(telegramBotToken, telegramWebhookAddress)
	if err != nil {
		return WireApplication{}, err
	}
	client := providers.ProvideYoutubeClient()
	factory := providers.LoggerFactoryProvider()
	youtubeService := services.NewYoutubeService(client, factory)
	telegramService := services.NewTelegramService(botAPI, youtubeService, factory)
	wireApplication := provideWireApplication(telegramService)
	return wireApplication, nil
}

// wire.go:

type WireApplication struct {
	telegramService *services.TelegramService
}

func provideWireApplication(telegramService *services.TelegramService) WireApplication {
	return WireApplication{telegramService: telegramService}
}

func (w WireApplication) Start() {
	w.telegramService.StartBot()
}
