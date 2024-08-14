package telegram

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog"
	"steplems-bot/lib"
	"steplems-bot/services/instagram"
	"strings"

	tbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"steplems-bot/services/youtube"
)

const EnableIgService = false

type TelegramService struct {
	api       *tbot.BotAPI
	ytService *youtube.YoutubeService
	igService *instagram.InstagramService
	logger    zerolog.Logger
	commands  map[string]TelegramCommand
}

func NewTelegramService(api *tbot.BotAPI,
	ytService *youtube.YoutubeService,
	igService *instagram.InstagramService,
	logger zerolog.Logger,
	cm *CommandMap) *TelegramService {
	return &TelegramService{api: api,
		ytService: ytService,
		igService: igService,
		logger:    logger,
		commands:  cm.commands,
	}
}

func (t *TelegramService) StartBot(ctx context.Context) error {
	uc := tbot.NewUpdate(0)
	updates := t.api.GetUpdatesChan(uc)
	if err := t.setCommands(); err != nil {
		t.logger.Error().Err(err).Msg("Failed to set commands")
		return err
	}
	if EnableIgService {
		go func() {
			restartChan := make(chan struct{}, 1)
			go func() {
				for range restartChan {
					go t.igService.Run("steplems", t.api, -1001373947640, restartChan)
				}
			}()
			restartChan <- struct{}{}
		}()
	}

	for update := range updates {
		update := update
		go func() {
			ctx, _ := context.WithCancel(ctx)
			err := t.OnUpdate(ctx, update)
			if err != nil {
				t.logger.Error().Err(err).Msg("Received error OnUpdate")
			}
		}()
	}
	return nil
}

func (t *TelegramService) OnUpdate(ctx context.Context, update tbot.Update) error {
	t.logger.Debug().Interface("chat", update.FromChat()).Interface("update", update).Msg("received update")
	if update.Message != nil {
		t.logger.Debug().Interface("from", update.SentFrom()).Str("message", update.Message.Text).Send()
	}
	if update.Message == nil {
		return nil
	}
	if t.ytService.IsYoutubeMessage(update) {
		c, err := t.ytService.MessageUpdate(update.Message)
		if err != nil {
			t.logger.Error().Err(err).Msg("Failed MessageUpdate")
			msg := tbot.NewMessage(update.FromChat().ID, fmt.Sprintf("youtube service error: %q", err.Error()))
			msg.ReplyToMessageID = update.Message.MessageID
			if _, err := t.api.Send(msg); err != nil {
				t.logger.Error().Err(err).Msg("failed to send")
			}
		} else {
			if _, err := t.api.Send(c); err != nil {
				t.logger.Error().Err(err).Msg("failed to send")
			}
		}
	}
	if update.Message.IsCommand() {
		command, ok := t.commands[update.Message.Command()]
		if !ok {
			return nil
		}
		cc := lib.NewChatContext(ctx, t, update)
		if err := command.Run(cc); err != nil {
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
