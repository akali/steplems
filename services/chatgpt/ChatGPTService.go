package chatgpt

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"steplems-bot/types"
)

type ChatGPTService struct {
	client          *openai.Client
	deepInfraClient *types.DeepInfraOpenAIClient
	logger          zerolog.Logger
}

func New(client *openai.Client, deepInfraClient *types.DeepInfraOpenAIClient, logger zerolog.Logger) *ChatGPTService {
	return &ChatGPTService{
		client:          client,
		deepInfraClient: deepInfraClient,
		logger:          logger,
	}
}

func (c *ChatGPTService) deepInfraAnswer(ctx context.Context, thread []openai.ChatCompletionMessage, model types.ModelStorage, config Config) (string, error) {
	client := (*openai.Client)(c.deepInfraClient)
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:     model.Model,
		Messages:  thread,
		MaxTokens: config.tokenSize,
	})

	if err != nil {
		return err.Error(), err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("got empty choices from response")
	}

	c.logger.Debug().Str("response", resp.Choices[0].Message.Content).Send()

	return resp.Choices[0].Message.Content, nil
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
	log.Debug().Interface("thread", thread).Msg("answering")
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
