package chatgpt

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
)

type ChatGPTService struct {
	client *openai.Client
}

func New(client *openai.Client) *ChatGPTService {
	return &ChatGPTService{
		client: client,
	}
}

func (c *ChatGPTService) Answer(ctx context.Context, question string) (string, error) {
	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{{
			Role:    openai.ChatMessageRoleAssistant,
			Content: question}},
	})

	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("got empty choices from response")
	}

	return resp.Choices[0].Message.Content, nil
}
