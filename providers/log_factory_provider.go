package providers

import (
	"github.com/google/wire"
	"github.com/olehan/kek"
	"github.com/olehan/kek/formatters/sugared"
	"os"
)

func LoggerFactoryProvider() *kek.Factory {
	return kek.NewFactory(os.Stdout, sugared.Formatter)
}

var LoggerFactoryProviderSet = wire.NewSet(LoggerFactoryProvider)
