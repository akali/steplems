package providers

import (
	"github.com/google/wire"
	"steplems-bot/types"

	"github.com/tmc/langchaingo/llms/openai"
)

func ProvideLangChainOpenAILLM(token types.OpenAIToken) (*openai.LLM, error) {
	return openai.New(openai.WithToken(string(token)))
}

func ProvideLangChainDeepInfraLLM(token types.DeepInfraToken) (*types.DeepInfraLangChainOpenAIClient, error) {
	client, err := openai.New(openai.WithToken(string(token)), openai.WithBaseURL("https://api.deepinfra.com/v1/openai"))
	if err != nil {
		return nil, err
	}
	return (*types.DeepInfraLangChainOpenAIClient)(client), nil
}

var LangChainWireSet = wire.NewSet(ProvideLangChainOpenAILLM, ProvideLangChainDeepInfraLLM)
