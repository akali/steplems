package services

import (
	"github.com/google/wire"
	"steplems-bot/services/chatgpt"
	"steplems-bot/services/spotify"
	"steplems-bot/services/telegram"
	"steplems-bot/services/youtube"
)

var ServicesSet = wire.NewSet(
	spotify.SpotifyServiceProviderSet,
	youtube.YoutubeServiceProvider,
	telegram.TelegramServiceSet,
	chatgpt.ChatGPTServiceProviderSet)
