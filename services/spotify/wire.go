package spotify

import "github.com/google/wire"

var SpotifyServiceProviderSet = wire.NewSet(NewSpotifyService, NewSpotifyAuthService)
