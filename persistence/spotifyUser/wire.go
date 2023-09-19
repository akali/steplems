package spotify_user

import "github.com/google/wire"

var SpotifyUserRepositoryProviderSet = wire.NewSet(NewSpotifyUserRepository)
