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

func NewChatContext(ctx context.Context, sender types.Sender, update tbot.Update, bot *tbot.BotAPI) *ChatContext {
	return &ChatContext{
		Ctx:    ctx,
		Sender: sender,
		bot:    bot,
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
