package types

import (
	"github.com/sashabaranov/go-openai"
	langchainopenai "github.com/tmc/langchaingo/llms/openai"
)

type DeepInfraToken string

type DeepInfraOpenAIClient openai.Client
type DeepInfraLangChainOpenAIClient langchainopenai.LLM

type ImGenModel struct {
	Model string
}
