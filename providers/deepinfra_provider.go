package providers

import (
	"github.com/google/wire"
	"github.com/sashabaranov/go-openai"
	"steplems-bot/lib/deepinfra"
	"steplems-bot/types"
)

func ProvideDeepInfraToken() (types.DeepInfraToken, error) {
	return ProvideEnvironmentVariable[types.DeepInfraToken]("DEEP_INFRA_TOKEN")()
}

var ProvideDeepInfraClient = deepinfra.NewClient

func ProvideDeepInfraOpenAIClient(token types.DeepInfraToken) *types.DeepInfraOpenAIClient {
	config := openai.DefaultConfig(string(token))
	config.BaseURL = "https://api.deepinfra.com/v1/openai"
	oConfig := openai.NewClientWithConfig(config)
	client := (*types.DeepInfraOpenAIClient)(oConfig)

	return client
}

var DeepInfraProviders = wire.NewSet(ProvideDeepInfraToken, ProvideDeepInfraClient, ProvideDeepInfraOpenAIClient)
