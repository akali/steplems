package commands

import "github.com/google/wire"

var CommandsProvider = wire.NewSet(
	NewHelpCommand,
	NewAuthorizeSpotifyCommand,
	NewNowPlayingCommand,
	NewChatGPTCommand,
	NewSetModelCommand,
	NewImGenCommand,
	NewTranscribeCommand,
	NewReasonCommand)
