package lib

import (
	"context"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hashicorp/go-multierror"
	"steplems-bot/types"
	"strings"
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

func (c *ChatContext) RespondImageURL(url string) {
	file := tbot.FileURL(url)
	msg := tbot.NewPhoto(c.Update.FromChat().ID, file)
	msg.ReplyToMessageID = c.Update.Message.MessageID
	_, err := c.Sender.Send(msg)
	c.Err = multierror.Append(c.Err, err)
}

func (cc *ChatContext) Text() string {
	input := cc.Update.Message.Text
	if strings.HasPrefix(input, "/") {
		_, after, found := strings.Cut(input, " ")
		if !found {
			return ""
		}
		return after
	}
	return input
}
