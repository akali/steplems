package commands

import (
	tbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"steplems-bot/services/spotify"
	"steplems-bot/types"
)

type AuthorizeSpotifyCommand struct {
	service *spotify.SpotifyService
}

func (h *AuthorizeSpotifyCommand) Run(service types.Sender, update tbot.Update) error {
	user, err := h.service.AuthorizeUser(update.Message.From.UserName)
	if err != nil {
		return err
	}
	msg := tbot.NewMessage(update.Message.Chat.ID, user.Username)
	_, err = service.Send(msg)
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
