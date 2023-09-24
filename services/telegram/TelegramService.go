package telegram

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"strings"

	"github.com/olehan/kek"

	tbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"steplems-bot/services/youtube"
)

type TelegramService struct {
	api       *tbot.BotAPI
	ytService *youtube.YoutubeService
	logger    *kek.Logger
	commands  map[string]TelegramCommand
}

func NewTelegramService(api *tbot.BotAPI,
	ytService *youtube.YoutubeService,
	kekFactory *kek.Factory,
	cm *CommandMap) *TelegramService {
	return &TelegramService{api: api,
		ytService: ytService,
		logger:    kekFactory.NewLogger("TelegramService"),
		commands:  cm.commands,
	}
}

func (t *TelegramService) StartBot(ctx context.Context) error {
	uc := tbot.NewUpdate(0)
	updates := t.api.GetUpdatesChan(uc)
	if err := t.setCommands(); err != nil {
		t.logger.Error.Println("Failed to set commands: ", err.Error())
		return err
	}
	for update := range updates {
		go func() {
			ctx, _ := context.WithCancel(ctx)
			err := t.OnUpdate(ctx, update)
			if err != nil {
				t.logger.Error.Println("Received error OnUpdate: ", err)
			}
		}()
	}
	return nil
}

func (t *TelegramService) OnUpdate(ctx context.Context, update tbot.Update) error {
	t.logger.Debug.Println("received an update: ", update)
	if update.Message == nil {
		return nil
	}
	if t.ytService.IsYoutubeMessage(update) {
		c, err := t.ytService.MessageUpdate(update.Message)
		if err != nil {
			t.logger.Error.Println("Failed MessageUpdate: ", err)
			msg := tbot.NewMessage(update.FromChat().ID, fmt.Sprintf("youtube service error: %q", err.Error()))
			msg.ReplyToMessageID = update.Message.MessageID
			if _, err := t.api.Send(msg); err != nil {
				t.logger.Error.Println("failed to send: ", err.Error())
			}
		} else {
			if _, err := t.api.Send(c); err != nil {
				t.logger.Error.Println("failed to send: ", err.Error())
			}
		}
	}
	if update.Message.IsCommand() {
		if err := t.commands[update.Message.Command()].Run(ctx, t, update); err != nil {
			msg := tbot.NewMessage(update.FromChat().ID, fmt.Sprintf("command error: %q", err.Error()))
			msg.ReplyToMessageID = update.Message.MessageID
			_, newErr := t.Send(msg)
			if newErr != nil {
				return multierror.Append(err, newErr)
			}
			return err
		}
	}
	return nil
}

func (t *TelegramService) setCommands() error {
	var cmds []tbot.BotCommand
	for _, command := range t.commands {
		cmd := ToAPITelegramCommand(command)
		cmd.Command = strings.TrimPrefix(strings.Split(cmd.Command, " ")[0], "/")
		cmds = append(cmds, cmd)
	}
	cfg := tbot.NewSetMyCommands(cmds...)
	_, err := t.api.Request(cfg)
	return err
}

func (t *TelegramService) Send(c tbot.Chattable) (tbot.Message, error) {
	return t.api.Send(c)
}
