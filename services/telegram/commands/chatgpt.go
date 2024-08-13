package commands

import (
	"context"
	"encoding/json"
	"fmt"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-openai"
	"steplems-bot/services/chatgpt"
	"steplems-bot/types"
	"strings"
)

var DefaultModel = types.ModelStorage{Model: openai.GPT3Dot5Turbo, Backend: "openai"}
var DeepInfraDefaultModel = types.ModelStorage{Model: "meta-llama/Meta-Llama-3.1-405B-Instruct", Backend: "deepinfra"}
var model = DefaultModel

type ChatGPTCommand struct {
	service *chatgpt.ChatGPTService
}

func NewChatGPTCommand(service *chatgpt.ChatGPTService) *ChatGPTCommand {
	return &ChatGPTCommand{
		service: service,
	}
}

func (c *ChatGPTCommand) Run(ctx context.Context, sender types.Sender, update tbot.Update) error {
	if update.Message == nil {
		return fmt.Errorf("this is not a message")
	}

	question := strings.TrimPrefix(update.Message.Text, "/"+c.Command())
	question = strings.TrimSuffix(strings.TrimPrefix(question, " "), " ")

	answer, err := c.service.Answer(ctx, question, model)
	if err != nil {
		return err
	}
	_, err = sender.Send(tbot.NewMessage(update.FromChat().ID, answer))
	return err
}

func (c *ChatGPTCommand) Command() string {
	return "chatgpt"
}

func (c *ChatGPTCommand) Description() string {
	return "Ask ChatGPT, anything!"
}

type SetModelCommand struct{}

func NewSetModelCommand() *SetModelCommand { return &SetModelCommand{} }

func (c *SetModelCommand) Command() string {
	return "setmodel"
}

func (c *SetModelCommand) Description() string {
	return "Set model and backend in json format. Like: `{'model': 'gpt-3.5-turbo', 'backend': 'openai'}` or `{'model': 'meta-llama/Meta-Llama-3.1-405B-Instruct', 'backend': 'deepinfra'}`"
}

func (c *SetModelCommand) Run(ctx context.Context, sender types.Sender, update tbot.Update) error {
	if update.Message == nil {
		return fmt.Errorf("this is not a message")
	}

	input := strings.TrimPrefix(update.Message.Text, "/"+c.Command())
	input = strings.TrimSuffix(strings.TrimPrefix(input, " "), " ")

	m := types.ModelStorage{}

	if input == "openai" {
		m = DefaultModel
	} else if input == "deepinfra" {
		m = DeepInfraDefaultModel
	} else {
		if err := json.Unmarshal([]byte(input), m); err != nil {
			return err
		}

		if m.Backend != "openai" && m.Backend != "deepinfra" {
			return fmt.Errorf("invalid backend format, must by `openai` or `deepinfra`, got %q", m.Backend)
		}
	}

	model = m
	_, err := sender.Send(tbot.NewMessage(update.FromChat().ID, fmt.Sprintf("model updated to: %q", model)))
	return err
}
