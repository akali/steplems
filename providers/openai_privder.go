package providers

import (
	"github.com/google/wire"
	"github.com/sashabaranov/go-openai"
	"steplems-bot/types"
)

func ProvideOpenAIToken() (types.OpenAIToken, error) {
	return ProvideEnvironmentVariable[types.OpenAIToken]("OPENAI_TOKEN")()
}

func ProvideOpenAIClient(token types.OpenAIToken) *openai.Client {
	return openai.NewClient(string(token))
}

var OpenAIProviders = wire.NewSet(ProvideOpenAIToken, ProvideOpenAIClient)
