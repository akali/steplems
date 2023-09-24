package spotify

import "github.com/google/wire"

var PersistenceSet = wire.NewSet(NewSpotifyUserRepository)
