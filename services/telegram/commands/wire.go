package commands

import "github.com/google/wire"

var CommandsProvider = wire.NewSet(NewHelpCommand, NewAuthorizeSpotifyCommand)
