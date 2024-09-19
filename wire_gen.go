// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"fmt"
	"log"
	"steplems-bot/lib/deepinfra"
	"steplems-bot/persistence/spotify_persistence"
	"steplems-bot/persistence/telegram_persistence"
	"steplems-bot/providers"
	"steplems-bot/services/chatgpt"
	deepinfra2 "steplems-bot/services/deepinfra"
	"steplems-bot/services/instagram"
	"steplems-bot/services/spotify"
	"steplems-bot/services/telegram"
	"steplems-bot/services/telegram/commands"
	"steplems-bot/services/youtube"
	"steplems-bot/types"
)

// Injectors from wire.go:

func NewWireApplication() (WireApplication, error) {
	port, err := providers.ProvidePort()
	if err != nil {
		return WireApplication{}, err
	}
	spotifyClientID, err := providers.ProvideSpotifyClientID()
	if err != nil {
		return WireApplication{}, err
	}
	spotifyClientSecret, err := providers.ProvideSpotifyClientSecret()
	if err != nil {
		return WireApplication{}, err
	}
	hostname, err := providers.ProvideHostname()
	if err != nil {
		return WireApplication{}, err
	}
	authenticator := providers.ProvideSpotifyAuth(spotifyClientID, spotifyClientSecret, hostname, port)
	consoleWriter := providers.LoggerOutputProvider()
	logger := providers.LoggerProvider(consoleWriter)
	spotifyAuthService := spotify.NewSpotifyAuthService(port, authenticator, logger)
	databaseConnectionURL, err := providers.ProvideDatabaseConnectionURL()
	if err != nil {
		return WireApplication{}, err
	}
	db, err := providers.ProvideDatabase(databaseConnectionURL, logger)
	if err != nil {
		return WireApplication{}, err
	}
	userRepository := spotify_persistence.NewSpotifyUserRepository(db)
	telegram_persistenceUserRepository := telegram_persistence.NewUserRepository(db)
	spotifyService := spotify.NewSpotifyService(port, spotifyAuthService, userRepository, telegram_persistenceUserRepository, authenticator, logger)
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
	youtubeService := youtube.NewYoutubeService(client, logger)
	goInstaConfigPath, err := providers.ProvideGoInstaConfigPath()
	if err != nil {
		return WireApplication{}, err
	}
	goinstaInstagram, err := providers.ProvideGoInsta(goInstaConfigPath)
	if err != nil {
		return WireApplication{}, err
	}
	instaCachePath, err := providers.ProvideInstaCachePath()
	if err != nil {
		return WireApplication{}, err
	}
	instagramService := instagram.New(goinstaInstagram, goInstaConfigPath, instaCachePath)
	authorizeSpotifyCommand := commands.NewAuthorizeSpotifyCommand(spotifyService)
	helpCommand := commands.NewHelpCommand()
	nowPlayingCommand := commands.NewNowPlayingCommand(spotifyService)
	openAIToken, err := providers.ProvideOpenAIToken()
	if err != nil {
		return WireApplication{}, err
	}
	openaiClient := providers.ProvideOpenAIClient(openAIToken)
	deepInfraToken, err := providers.ProvideDeepInfraToken()
	if err != nil {
		return WireApplication{}, err
	}
	deepInfraOpenAIClient := providers.ProvideDeepInfraOpenAIClient(deepInfraToken)
	deepInfraLangChainOpenAIClient, err := providers.ProvideLangChainDeepInfraLLM(deepInfraToken)
	if err != nil {
		return WireApplication{}, err
	}
	chatGPTService := chatgpt.New(openaiClient, deepInfraOpenAIClient, deepInfraLangChainOpenAIClient, logger)
	messageRepository := telegram_persistence.NewMessageRepository(db)
	chatGPTCommand := commands.NewChatGPTCommand(chatGPTService, messageRepository, telegram_persistenceUserRepository)
	setModelCommand := commands.NewSetModelCommand()
	deepinfraClient := deepinfra.NewClient(deepInfraToken, logger)
	deepInfraService := deepinfra2.NewStableDiffusionService(deepinfraClient, logger)
	imGenCommand := commands.NewImGenCommand(deepInfraService)
	transcribeCommand := commands.NewTranscribeCommand(deepInfraService)
	reasonService := chatgpt.NewReasonService(chatGPTService)
	reasonCommand := commands.NewReasonCommand(reasonService)
	commandMap := telegram.NewCommandMap(authorizeSpotifyCommand, helpCommand, nowPlayingCommand, chatGPTCommand, setModelCommand, imGenCommand, transcribeCommand, reasonCommand)
	telegramService := telegram.NewTelegramService(botAPI, youtubeService, instagramService, logger, commandMap)
	wireApplication := provideWireApplication(spotifyService, spotifyAuthService, telegramService, hostname, userRepository, telegram_persistenceUserRepository, messageRepository)
	return wireApplication, nil
}

// wire.go:

type WireApplication struct {
	telegramService *telegram.TelegramService
	sUserRepo       *spotify_persistence.UserRepository
	tUserRepo       *telegram_persistence.UserRepository
	messageRepo     *telegram_persistence.MessageRepository
	spotifyService  *spotify.SpotifyService
	hostname        types.Hostname
	authService     *spotify.SpotifyAuthService
}

func provideWireApplication(spotifyService *spotify.SpotifyService, authService *spotify.SpotifyAuthService, telegramService *telegram.TelegramService, hostname types.Hostname, sUserRepo *spotify_persistence.UserRepository, tUserRepo *telegram_persistence.UserRepository, messageRepository *telegram_persistence.MessageRepository) WireApplication {
	return WireApplication{
		spotifyService:  spotifyService,
		authService:     authService,
		telegramService: telegramService,
		sUserRepo:       sUserRepo,
		tUserRepo:       tUserRepo,
		messageRepo:     messageRepository,
		hostname:        hostname}
}

func (w WireApplication) Start(command2 string) error {
	ctx := context.Background()

	switch command2 {
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
