package providers

import (
	"steplems-bot/types"

	"github.com/google/wire"
)

func ProvideBotToken() (types.TelegramBotToken, error) {
	return ProvideEnvironmentVariable[types.TelegramBotToken]("TELEGRAM_TOKEN")()
}

func ProvideBotWebhook() (types.TelegramWebhookAddress, error) {
	return "", nil
}

var TelegramKeyProviderSet = wire.NewSet(ProvideBotToken, ProvideBotWebhook)
