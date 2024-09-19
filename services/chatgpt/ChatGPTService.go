package chatgpt

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/sashabaranov/go-openai"
	"github.com/tmc/langchaingo/llms"
	langchainopenai "github.com/tmc/langchaingo/llms/openai"
	"steplems-bot/types"
)

type ChatGPTService struct {
	client                         *openai.Client
	deepInfraClient                *types.DeepInfraOpenAIClient
	deepInfraLangChainOpenAIClient *types.DeepInfraLangChainOpenAIClient
	logger                         zerolog.Logger
}

func New(client *openai.Client, deepInfraClient *types.DeepInfraOpenAIClient, deepInfraLangChainOpenAIClient *types.DeepInfraLangChainOpenAIClient, logger zerolog.Logger) *ChatGPTService {
	return &ChatGPTService{client: client, deepInfraClient: deepInfraClient, deepInfraLangChainOpenAIClient: deepInfraLangChainOpenAIClient, logger: logger}
}

func openaiRoleToLCRole(role string) llms.ChatMessageType {
	switch role {
	case openai.ChatMessageRoleAssistant:
		return llms.ChatMessageTypeAI
	case openai.ChatMessageRoleUser:
		return llms.ChatMessageTypeHuman
	case openai.ChatMessageRoleSystem:
		return llms.ChatMessageTypeSystem
	case openai.ChatMessageRoleFunction:
		return llms.ChatMessageTypeFunction
	}
	return llms.ChatMessageTypeHuman
}

func (c *ChatGPTService) deepInfraAnswer(ctx context.Context, thread []openai.ChatCompletionMessage, model types.ModelStorage, config Config) (string, error) {
	client := (*langchainopenai.LLM)(c.deepInfraLangChainOpenAIClient)
	var messages []llms.MessageContent
	for _, message := range thread {
		messages = append(messages, llms.MessageContent{
			Role: openaiRoleToLCRole(message.Role),
			Parts: []llms.ContentPart{
				llms.TextContent{Text: message.Content},
			},
		})
	}
	resp, err := client.GenerateContent(ctx, messages, llms.WithJSONMode(), llms.WithMaxTokens(config.tokenSize), llms.WithModel(model.Model), llms.WithJSONMode())

	if err != nil {
		return err.Error(), err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("got empty choices from response")
	}

	return resp.Choices[0].Content, nil
}

type Config struct {
	tokenSize int
}

type Option func(*Config)

func buildConfig(opts ...Option) Config {
	config := Config{}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}

func WithMaxTokenSize(tokenSize int) Option {
	return func(c *Config) {
		c.tokenSize = tokenSize
	}
}

func (c *ChatGPTService) Answer(ctx context.Context, thread []openai.ChatCompletionMessage, model types.ModelStorage, opts ...Option) (string, error) {
	config := buildConfig(opts...)
	if model.Backend == "deepinfra" {
		return c.deepInfraAnswer(ctx, thread, model, config)
	}
	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: thread,
	})

	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("got empty choices from response")
	}

	return resp.Choices[0].Message.Content, nil
}
