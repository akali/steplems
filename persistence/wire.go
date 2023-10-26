package persistence

import (
	"github.com/google/wire"
	"steplems-bot/persistence/spotify_persistence"
	"steplems-bot/persistence/telegram_persistence"
)

var PersistenceSet = wire.NewSet(telegram_persistence.TelegramPersistenceSet, spotify_persistence.SpotifyPersistenceSet)
