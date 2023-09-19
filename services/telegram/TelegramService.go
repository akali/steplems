package services

import (
	"fmt"
	"steplems-bot/types"
	"strings"

	"github.com/google/wire"
	"github.com/olehan/kek"

	tbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramService struct {
	api       *tbot.BotAPI
	ytService *YoutubeService
	logger    *kek.Logger
	commands  []types.TelegramCommand
}

func NewTelegramService(api *tbot.BotAPI, ytService *YoutubeService, kekFactory *kek.Factory) *TelegramService {
	return &TelegramService{api: api,
		ytService: ytService,
		logger:    kekFactory.NewLogger("TelegramService")}
}

func (t *TelegramService) StartBot() {
	uc := tbot.NewUpdate(0)
	updates := t.api.GetUpdatesChan(uc)
	for update := range updates {
		go t.OnUpdate(update)
	}
}

func (t *TelegramService) RegisterCommand(commands ...types.TelegramCommand) error {
	for _, cmd := range commands {
		cmd.Command = strings.TrimPrefix(strings.Split(cmd.Command, " ")[0], "/")
		t.commands = append(t.commands, cmd)
	}

	return t.setCommands()
}

func (t *TelegramService) OnUpdate(update tbot.Update) {
	t.logger.Debug.Println("received an update: ", update)
	if t.ytService.IsYoutubeMessage(update) {
		c, err := t.ytService.MessageUpdate(update.Message)
		if err != nil {
			t.logger.Error.Println("Failed MessageUpdate: ", err)
			msg := tbot.NewMessage(update.Message.Chat.ID, fmt.Sprintf("youtube service error: %q", err.Error()))
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
}

func (t *TelegramService) setCommands() error {
	cmds := []tbot.BotCommand{}
	for _, command := range t.commands {
		cmds = append(cmds, command.ToAPITelegramCommand())
	}
	cfg := tbot.NewSetMyCommands(cmds...)
	_, err := t.api.Send(cfg)
	return err
}

var TelegramServiceProvider = wire.NewSet(NewTelegramService)
