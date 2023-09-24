package commands

import (
	"context"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"steplems-bot/services/spotify"
	"steplems-bot/types"
)

type AuthorizeSpotifyCommand struct {
	spotifyService *spotify.SpotifyService
}

func (h *AuthorizeSpotifyCommand) Run(ctx context.Context, sender types.Sender, update tbot.Update) error {
	if update.Message.Chat.Type != "private" {
		msg := tbot.NewMessage(update.FromChat().ID, "To authorize spotify, chat with bot in DM.")
		_, err := sender.Send(msg)
		return err
	}

	user, err := h.spotifyService.AuthorizeUser(ctx, sender, update)
	if err != nil {
		return err
	}
	msg := tbot.NewMessage(update.FromChat().ID, user.ID)
	_, err = sender.Send(msg)
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
