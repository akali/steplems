package providers

import (
	"steplems-bot/types"

	"github.com/google/wire"
)
import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ProvideTgBot(token types.TelegramBotToken, _ types.TelegramWebhookAddress) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(string(token))
	if err != nil {
		return nil, err
	}
	return bot, nil
}

var TelegramBotProviderSet = wire.NewSet(ProvideTgBot)
