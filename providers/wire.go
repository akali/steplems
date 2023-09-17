package providers

import "github.com/google/wire"

var ProvidersSet = wire.NewSet(
	TelegramBotProviderSet,
	TelegramKeyProviderSet,
	YoutubeClientProvider,
	LoggerFactoryProviderSet)
