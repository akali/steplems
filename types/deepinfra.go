package types

import "github.com/sashabaranov/go-openai"

type DeepInfraToken string

type DeepInfraOpenAIClient openai.Client

type ImGenModel struct {
	Model string
}
