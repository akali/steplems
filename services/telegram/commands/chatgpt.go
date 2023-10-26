package commands

import (
	"context"
	"fmt"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"steplems-bot/services/chatgpt"
	"steplems-bot/types"
)

type ChatGPTCommand struct {
	service *chatgpt.ChatGPTService
}

func NewChatGPTCommand(service *chatgpt.ChatGPTService) *ChatGPTCommand {
	return &ChatGPTCommand{
		service: service,
	}
}

func (c *ChatGPTCommand) Run(ctx context.Context, sender types.Sender, update tbot.Update) error {
	if update.Message == nil {
		return fmt.Errorf("this is not a message")
	}

	question := update.Message.Text

	answer, err := c.service.Answer(ctx, question)
	if err != nil {
		return err
	}
	_, err = sender.Send(tbot.NewMessage(update.FromChat().ID, answer))
	return err
}

func (c *ChatGPTCommand) Command() string {
	return "chatgpt"
}

func (c *ChatGPTCommand) Description() string {
	return "Ask ChatGPT, anything!"
}
