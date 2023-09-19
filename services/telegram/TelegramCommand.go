package telegram

import (
	tbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/wire"
	"steplems-bot/services/spotify"
)

type CommandMap struct {
	commands map[string]TelegramCommand
}

func NewCommandMap(authorizeSpotifyCommand *AuthorizeSpotifyCommand) *CommandMap {
	commands := []TelegramCommand{
		Help{},
		authorizeSpotifyCommand,
	}
	cm := CommandMap{commands: make(map[string]TelegramCommand)}
	for _, command := range commands {
		cm.commands[command.Command()] = command
	}
	return &cm
}

var CommandMapProvider = wire.NewSet(NewCommandMap, NewAuthorizeSpotifyCommand)

type Help struct{}

func (h Help) Run(service *TelegramService, update tbot.Update) error {
	msg := tbot.NewMessage(update.Message.Chat.ID, "help requested")
	_, err := service.api.Send(msg)
	return err
}

func (h Help) Command() string {
	return "help"
}

func (h Help) Description() string {
	return "Get help."
}

type AuthorizeSpotifyCommand struct {
	service *spotify.SpotifyService
}

func (h *AuthorizeSpotifyCommand) Run(service *TelegramService, update tbot.Update) error {
	user, err := h.service.AuthorizeUser(update.Message.From.UserName)
	if err != nil {
		return err
	}
	msg := tbot.NewMessage(update.Message.Chat.ID, user.Username)
	_, err = service.api.Send(msg)
	return err
}

func (c *AuthorizeSpotifyCommand) Command() string {
	return "authorize"
}

func (c *AuthorizeSpotifyCommand) Description() string {
	return "Authorize your spotify account."
}

func NewAuthorizeSpotifyCommand(service *spotify.SpotifyService) *AuthorizeSpotifyCommand {
	return &AuthorizeSpotifyCommand{service}
}

type TelegramCommand interface {
	Run(*TelegramService, tbot.Update) error
	Command() string
	Description() string
}

func ToAPITelegramCommand(cmd TelegramCommand) tbot.BotCommand {
	return tbot.BotCommand{
		Command:     cmd.Command(),
		Description: cmd.Description(),
	}
}
