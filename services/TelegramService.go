package services

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/wire"
	"github.com/olehan/kek"
)

type TelegramService struct {
	api       *tgbotapi.BotAPI
	ytService *YoutubeService
	logger    *kek.Logger
}

func NewTelegramService(api *tgbotapi.BotAPI, ytService *YoutubeService, kekFactory *kek.Factory) *TelegramService {
	return &TelegramService{api: api,
		ytService: ytService,
		logger:    kekFactory.NewLogger("TelegramService")}
}

func (t *TelegramService) StartBot() {
	uc := tgbotapi.NewUpdate(0)
	updates := t.api.GetUpdatesChan(uc)
	go func() {
		for update := range t.ytService.Updates() {
			t.logger.Info.Println("got update message from youtube")
			update(t.api)
		}
	}()
	for update := range updates {
		go t.OnUpdate(update)
	}
}

func (t *TelegramService) OnUpdate(update tgbotapi.Update) {
	t.logger.Debug.Println("received an update: ", update)
	if t.ytService.YoutubeMessage(update) {
		if err := t.ytService.MessageUpdate(update.Message); err != nil {
			t.logger.Error.Println("Failed MessageUpdate: ", err)
		}
	}
}

var TelegramServiceProvider = wire.NewSet(NewTelegramService)
