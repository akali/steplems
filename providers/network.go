package providers

import (
	"github.com/google/wire"
	"steplems-bot/types"
)

func ProvideWebhookAddress() (types.WebhookAddress, error) {
	return ProvideEnvironmentVariable[types.WebhookAddress]("WEBHOOK")()
}

var NetworkProviders = wire.NewSet(ProvideWebhookAddress)
