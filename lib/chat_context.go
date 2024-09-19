package lib

import (
	"bytes"
	"context"
	"fmt"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hashicorp/go-multierror"
	"io"
	"net/http"
	"steplems-bot/types"
	"strings"
)

type ChatContext struct {
	Sender types.Sender
	Update tbot.Update
	bot    *tbot.BotAPI

	Ctx context.Context
	Err *multierror.Error
}

func (cc *ChatContext) Error() error {
	return cc.Err.ErrorOrNil()
}

func NewChatContext(ctx context.Context, sender types.Sender, update tbot.Update, bot *tbot.BotAPI) *ChatContext {
	return &ChatContext{
		Ctx:    ctx,
		Sender: sender,
		bot:    bot,
		Update: update,
		Err:    nil,
	}
}

func (cc *ChatContext) send(c tbot.Chattable) tbot.Message {
	message, err := cc.Sender.Send(c)
	cc.Err = multierror.Append(cc.Err, err)
	return message
}

func (c *ChatContext) RespondText(message string) tbot.Message {
	return c.send(tbot.NewMessage(c.Update.FromChat().ID, message))
}

func (c *ChatContext) ReplyText(message string) tbot.Message {
	msg := tbot.NewMessage(c.Update.FromChat().ID, message)
	msg.ReplyToMessageID = c.Update.Message.MessageID
	return c.send(msg)
}

func (c *ChatContext) EditMessage(message tbot.Message, newText string) tbot.Message {
	edit := tbot.NewEditMessageText(message.Chat.ID, message.MessageID, newText)
	edit.ParseMode = tbot.ModeMarkdown
	return c.send(edit)
}

func (c *ChatContext) RespondImageURL(url string) {
	file := tbot.FileURL(url)
	msg := tbot.NewPhoto(c.Update.FromChat().ID, file)
	msg.ReplyToMessageID = c.Update.Message.MessageID
	c.send(msg)
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

func (cc *ChatContext) GetFile(fileId string) ([]byte, error) {
	tFile, err := cc.bot.GetFile(tbot.FileConfig{
		FileID: fileId,
	})
	if err != nil {
		return nil, err
	}
	downloadURL := tFile.Link(cc.bot.Token)
	req, err := http.NewRequest("GET", downloadURL, bytes.NewBufferString(""))
	if err != nil {
		return nil, err
	}
	resp, err := cc.bot.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to generate image: status code %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (cc *ChatContext) RespondImage(b []byte) {
	photoMessage := tbot.NewPhoto(cc.Update.FromChat().ID, tbot.FileBytes{
		Name:  "flux.png",
		Bytes: b,
	})
	photoMessage.ReplyToMessageID = cc.Update.Message.MessageID
	cc.send(photoMessage)
}
