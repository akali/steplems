package telegram

import "github.com/google/wire"

var PersistenceSet = wire.NewSet(NewUserRepository)
