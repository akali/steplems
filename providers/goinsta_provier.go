package providers

import (
	"github.com/Davincible/goinsta/v3"
	"github.com/google/wire"
	"steplems-bot/types"
)

func ProvideGoInstaConfigPath() (types.GoInstaConfigPath, error) {
	return ProvideEnvironmentVariable[types.GoInstaConfigPath]("GOINSTA_CONFIG_PATH")()
}

func ProvideGoInsta(path types.GoInstaConfigPath) (*goinsta.Instagram, error) {
	return goinsta.Import(string(path))
}

var GoInstaProviders = wire.NewSet(ProvideGoInstaConfigPath, ProvideGoInsta)
