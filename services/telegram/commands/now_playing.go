package commands

import (
	"steplems-bot/lib"
	"steplems-bot/services/spotify"
)

type NowPlayingCommand struct {
	service *spotify.SpotifyService
}

func NewNowPlayingCommand(service *spotify.SpotifyService) *NowPlayingCommand {
	return &NowPlayingCommand{service: service}
}

func (c *NowPlayingCommand) Run(cc *lib.ChatContext) error {
	return c.service.NowPlaying(cc.Ctx, cc.Sender, cc.Update)
}

func (c *NowPlayingCommand) Command() string {
	return "nowplaying"
}

func (c *NowPlayingCommand) Description() string {
	return "Now playing on Spotify."
}
