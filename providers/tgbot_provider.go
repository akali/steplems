package providers

import (
	"github.com/google/wire"
	"steplems-bot/types"
)
import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ProvideTgBot(token types.TelegramBotToken, _ types.TelegramWebhookAddress) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(string(token))
	if err != nil {
		return nil, err
	}
	//wh, err := tgbotapi.NewWebhook(string(webhookAddress) + bot.Token)
	//if err != nil {
	//	return nil, err
	//}
	//
	//_, err = bot.Request(wh)
	//if err != nil {
	//	return nil, err
	//}
	//
	//info, err := bot.GetWebhookInfo()
	//if err != nil {
	//	return nil, err
	//}
	//
	//if info.LastErrorDate != 0 {
	//	return nil, fmt.Errorf("telegram callback failed: %s", info.LastErrorMessage)
	//}
	//
	return bot, nil
}

var TelegramBotProviderSet = wire.NewSet(ProvideTgBot)
