package providers

import (
	"github.com/google/wire"
	"steplems-bot/types"
)

func ProvideWebhookAddress() (types.WebhookAddress, error) {
	return ProvideEnvironmentVariable[types.WebhookAddress]("WEBHOOK")()
}

func ProvideHostname() (types.Hostname, error) {
	return ProvideEnvironmentVariable[types.Hostname]("HOSTNAME")()
}

func ProvidePort() (types.Port, error) {
	return ProvideEnvironmentVariable[types.Port]("PORT")()
}

var NetworkProviders = wire.NewSet(ProvideWebhookAddress, ProvideHostname, ProvidePort)
