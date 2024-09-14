//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"log"
	"steplems-bot/persistence"
	"steplems-bot/persistence/spotify_persistence"
	telegram2 "steplems-bot/persistence/telegram_persistence"
	"steplems-bot/providers"
	"steplems-bot/services"
	spotify2 "steplems-bot/services/spotify"
	"steplems-bot/services/telegram"
	"steplems-bot/types"
)

type WireApplication struct {
	telegramService *telegram.TelegramService
	sUserRepo       *spotify_persistence.UserRepository
	tUserRepo       *telegram2.UserRepository
	messageRepo     *telegram2.MessageRepository
	spotifyService  *spotify2.SpotifyService
	hostname        types.Hostname
	authService     *spotify2.SpotifyAuthService
}

func NewWireApplication() (WireApplication, error) {
	wire.Build(provideWireApplication,
		persistence.PersistenceSet,
		providers.ProvidersSet,
		services.ServicesSet)
	return WireApplication{}, nil
}

func provideWireApplication(spotifyService *spotify2.SpotifyService, authService *spotify2.SpotifyAuthService, telegramService *telegram.TelegramService, hostname types.Hostname, sUserRepo *spotify_persistence.UserRepository, tUserRepo *telegram2.UserRepository, messageRepository *telegram2.MessageRepository) WireApplication {
	return WireApplication{
		spotifyService:  spotifyService,
		authService:     authService,
		telegramService: telegramService,
		sUserRepo:       sUserRepo,
		tUserRepo:       tUserRepo,
		messageRepo:     messageRepository,
		hostname:        hostname}
}

func (w WireApplication) Start(command string) error {
	ctx := context.Background()

	switch command {
	case "runbot":
		return w.runbot(ctx)
	case "printEmails":
		return w.printEmails(ctx)
	case "migrate":
		return w.migrate()
	}

	return nil
}

func (w WireApplication) migrate() error {
	migratables := []types.MigrationRunner{
		w.sUserRepo,
		w.tUserRepo,
		w.messageRepo,
	}

	for _, m := range migratables {
		if err := m.RunMigrations(); err != nil {
			return err
		}
	}
	return nil
}

func (w WireApplication) runbot(ctx context.Context) error {
	log.Printf("Starting application with hostname=%s\n", string(w.hostname))
	go w.authService.Serve()

	return w.telegramService.StartBot(ctx)
}

func (w WireApplication) printEmails(ctx context.Context) error {
	users := w.sUserRepo.FindAll()
	for _, user := range users {
		client, err := w.spotifyService.CreateClient(ctx, user)
		if err != nil {
			return err
		}
		puser, err := client.CurrentUser(ctx)
		if err != nil {
			return err
		}
		fmt.Println(puser.Email)
	}
	return nil
}
