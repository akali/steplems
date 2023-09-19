package telegram

import "github.com/google/wire"

var TelegramServiceSet = wire.NewSet(CommandMapProvider, NewTelegramService)
