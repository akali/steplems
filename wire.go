//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"log"
	"steplems-bot/persistence"
	"steplems-bot/persistence/spotify"
	telegram2 "steplems-bot/persistence/telegram"
	"steplems-bot/services"
	spotify2 "steplems-bot/services/spotify"
	"steplems-bot/services/telegram"
	"steplems-bot/types"
)
import "steplems-bot/providers"

type WireApplication struct {
	telegramService *telegram.TelegramService
	sUserRepo       *spotify.UserRepository
	tUserRepo       *telegram2.UserRepository
	hostname        types.Hostname
	authService     *spotify2.SpotifyAuthService
}

func provideWireApplication(authService *spotify2.SpotifyAuthService, telegramService *telegram.TelegramService, hostname types.Hostname, sUserRepo *spotify.UserRepository, tUserRepo *telegram2.UserRepository) WireApplication {
	return WireApplication{authService: authService, telegramService: telegramService, sUserRepo: sUserRepo, tUserRepo: tUserRepo, hostname: hostname}
}

func NewWireApplication() (WireApplication, error) {
	wire.Build(provideWireApplication,
		providers.ProvidersSet,
		services.ServicesSet,
		persistence.PersistenceSet)
	return WireApplication{}, nil
}

func (w WireApplication) Start() error {
	ctx := context.Background()

	migratables := []types.MigrationRunner{
		w.sUserRepo,
		w.tUserRepo,
	}

	for _, m := range migratables {
		if err := m.RunMigrations(); err != nil {
			return err
		}
	}

	log.Printf("Starting application with hostname=%s\n", string(w.hostname))

	go w.authService.Serve()

	return w.telegramService.StartBot(ctx)
}
