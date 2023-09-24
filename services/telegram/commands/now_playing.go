package commands

import (
	"context"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"steplems-bot/services/spotify"
	"steplems-bot/types"
)

type NowPlayingCommand struct {
	service *spotify.SpotifyService
}

func NewNowPlayingCommand(service *spotify.SpotifyService) *NowPlayingCommand {
	return &NowPlayingCommand{service: service}
}

func (c *NowPlayingCommand) Run(ctx context.Context, sender types.Sender, update tbot.Update) error {
	return c.service.NowPlaying(ctx, sender, update)
}

func (c *NowPlayingCommand) Command() string {
	return "nowplaying"
}

func (c *NowPlayingCommand) Description() string {
	return "Now playing on Spotify."
}
