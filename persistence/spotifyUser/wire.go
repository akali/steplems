package spotifyUser

import "github.com/google/wire"

var SpotifyUserRepositoryProviderSet = wire.NewSet(NewSpotifyUserRepository)
