package lib

import (
	"context"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hashicorp/go-multierror"
	"steplems-bot/types"
)

type ChatContext struct {
	Sender types.Sender
	Update tbot.Update

	Ctx context.Context
	Err *multierror.Error
}

func NewChatContext(ctx context.Context, sender types.Sender, update tbot.Update) *ChatContext {
	return &ChatContext{
		Ctx:    ctx,
		Sender: sender,
		Update: update,
		Err:    nil,
	}
}

func (c *ChatContext) RespondText(message string) {
	_, err := c.Sender.Send(tbot.NewMessage(c.Update.FromChat().ID, message))
	c.Err = multierror.Append(c.Err, err)
}
