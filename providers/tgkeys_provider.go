package providers

import (
	"fmt"
	"github.com/google/wire"
	"os"
	"steplems-bot/types"
)

func ProvideBotToken() (types.TelegramBotToken, error) {
	if token, ok := os.LookupEnv("TELEGRAM_TOKEN"); ok {
		return types.TelegramBotToken(token), nil
	} else {
		return "", fmt.Errorf("token not found in environment variables")
	}
}

func ProvideBotWebhook() (types.TelegramWebhookAddress, error) {
	return "", nil
}

var TelegramKeyProviderSet = wire.NewSet(ProvideBotToken, ProvideBotWebhook)
