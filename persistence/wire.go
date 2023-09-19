//go:build wireinject
// +build wireinject

package persistence

import (
	"github.com/google/wire"
	"steplems-bot/persistence/spotifyUser"
)

var PersistenceSet = wire.NewSet(spotifyUser.SpotifyUserRepositoryProviderSet)
