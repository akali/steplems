package commands

import (
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"steplems-bot/lib"
	"steplems-bot/services/chatgpt"
	"steplems-bot/types"
	"strings"
)

var DefaultModel = types.ModelStorage{Model: openai.GPT3Dot5Turbo, Backend: "openai", ImGenModel: "stablediffusion"}
var DeepInfraDefaultModel = types.ModelStorage{Model: "meta-llama/Meta-Llama-3.1-405B-Instruct", Backend: "deepinfra", ImGenModel: "stablediffusion"}
var model = DefaultModel

var (
	imGenSD   = "stablediffusion"
	imGenFlux = "flux"
)

type ChatGPTCommand struct {
	service *chatgpt.ChatGPTService
}

func NewChatGPTCommand(service *chatgpt.ChatGPTService) *ChatGPTCommand {
	return &ChatGPTCommand{
		service: service,
	}
}

func (c *ChatGPTCommand) Run(cc *lib.ChatContext) error {
	if cc.Update.Message == nil {
		return fmt.Errorf("this is not a message")
	}

	question := strings.TrimPrefix(cc.Update.Message.Text, "/"+c.Command())
	question = strings.TrimSuffix(strings.TrimPrefix(question, " "), " ")

	answer, err := c.service.Answer(cc.Ctx, question, model)
	if err != nil {
		return err
	}
	cc.RespondText(answer)
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
	return "Set model and backend in json format. Like: `{'model': 'gpt-3.5-turbo', 'backend': 'openai', 'imgenmodel': 'stablediffusion'}`"
}

func (c *SetModelCommand) Run(cc *lib.ChatContext) error {
	if cc.Update.Message == nil {
		return fmt.Errorf("this is not a message")
	}

	input := strings.TrimPrefix(cc.Update.Message.Text, "/"+c.Command())
	input = strings.TrimSuffix(strings.TrimPrefix(input, " "), " ")

	m := &types.ModelStorage{}

	switch input {
	case "openai":
		oldImGen := model.ImGenModel
		m = &DefaultModel
		m.ImGenModel = oldImGen
	case "deepinfra":
		oldImGen := model.ImGenModel
		m = &DeepInfraDefaultModel
		m.ImGenModel = oldImGen
	case imGenSD:
		m = &model
		m.ImGenModel = imGenSD
	case imGenFlux:
		m = &model
		m.ImGenModel = imGenSD
	default:
		if err := json.Unmarshal([]byte(input), m); err != nil {
			return err
		}

		if m.Backend != "openai" && m.Backend != "deepinfra" {
			return fmt.Errorf("invalid backend format, must be `openai` or `deepinfra`, got %q", m.Backend)
		}
		if m.ImGenModel != imGenFlux && m.ImGenModel != imGenSD {
			return fmt.Errorf("invalid imgen model format, must be `stablediffusion` or `flux`, got %q", m.ImGenModel)
		}
	}

	model = *m
	cc.RespondText(fmt.Sprintf("model updated to: %q", model))
	return nil
}
