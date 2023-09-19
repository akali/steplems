package types

import (
	tbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBotToken string
type TelegramWebhookAddress string
type TelegramBotAction func(api *tbot.BotAPI)
type Sender interface {
	Send(tbot.Chattable) (tbot.Message, error)
}
