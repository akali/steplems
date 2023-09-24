package commands

import (
	"context"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"steplems-bot/types"
)

type HelpCommand struct{}

func (h *HelpCommand) Run(ctx context.Context, service types.Sender, update tbot.Update) error {
	msg := tbot.NewMessage(update.FromChat().ID, "help requested")
	_, err := service.Send(msg)
	return err
}

func (h *HelpCommand) Command() string {
	return "help"
}

func (h *HelpCommand) Description() string {
	return "Get help."
}

func NewHelpCommand() *HelpCommand {
	return &HelpCommand{}
}
