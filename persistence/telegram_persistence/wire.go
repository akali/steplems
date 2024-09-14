package telegram_persistence

import "github.com/google/wire"

var TelegramPersistenceSet = wire.NewSet(NewUserRepository, NewMessageRepository)
