//go:build wireinject
// +build wireinject

package persistence

import (
	"github.com/google/wire"
	"steplems-bot/persistence/spotify"
	"steplems-bot/persistence/telegram"
)

var PersistenceSet = wire.NewSet(telegram.PersistenceSet, spotify.PersistenceSet)
