package commands

import (
	"encoding/json"
	"fmt"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"slices"
	"steplems-bot/lib"
	"steplems-bot/persistence/telegram_persistence"
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
	repo    *telegram_persistence.MessageRepository
	uRepo   *telegram_persistence.UserRepository
}

func NewChatGPTCommand(service *chatgpt.ChatGPTService, repo *telegram_persistence.MessageRepository, uRepo *telegram_persistence.UserRepository) *ChatGPTCommand {
	return &ChatGPTCommand{
		service: service,
		repo:    repo,
		uRepo:   uRepo,
	}
}

func (c *ChatGPTCommand) Run(cc *lib.ChatContext) error {
	if cc.Update.Message == nil {
		return fmt.Errorf("this is not a message")
	}

	if !cc.Update.Message.IsCommand() {
		return c.Reply(cc)
	}

	message := telegram_persistence.FromTelegramMessage(cc.Update.Message)
	if err := c.repo.Create(message); err != nil {
		return fmt.Errorf("failed to save message in database: %w", err)
	}

	question := strings.TrimPrefix(cc.Update.Message.Text, "/"+c.Command())
	question = strings.TrimSuffix(strings.TrimPrefix(question, " "), " ")

	answer, err := c.service.Answer(cc.Ctx, []openai.ChatCompletionMessage{{
		Role:    openai.ChatMessageRoleUser,
		Content: question,
	}}, model)
	if err != nil {
		return err
	}
	replyMessage := cc.ReplyText(answer)
	telegramReplyMessage := telegram_persistence.FromTelegramMessage(&replyMessage)
	log.Debug().Interface("reply_message", replyMessage).Msg("reply")
	_, err = c.uRepo.GetOrCreate(telegramReplyMessage.FromUserID.Int64, telegram_persistence.FromExternalTelegramUser(replyMessage.From, replyMessage.Chat))
	if err != nil {
		log.Err(err).Send()
		return err
	}
	if err := c.repo.Create(telegramReplyMessage); err != nil {
		return err
	}
	return err
}

func fromThread(thread []telegram_persistence.Message) []openai.ChatCompletionMessage {
	var result []openai.ChatCompletionMessage
	for _, message := range thread {
		chatCompletionMessage := openai.ChatCompletionMessage{
			Content: message.Text,
		}
		if message.FromUser.IsBot {
			chatCompletionMessage.Role = openai.ChatMessageRoleAssistant
		} else {
			chatCompletionMessage.Role = openai.ChatMessageRoleUser
		}
		result = append(result, chatCompletionMessage)
	}
	slices.Reverse(result)
	return result
}

func (c *ChatGPTCommand) Reply(cc *lib.ChatContext) error {
	if cc.Update.Message == nil {
		return fmt.Errorf("this is not a message")
	}

	if cc.Update.Message.ReplyToMessage == nil {
		return fmt.Errorf("no reply in message")
	}

	message := telegram_persistence.FromTelegramMessage(cc.Update.Message)
	if err := c.repo.Create(message); err != nil {
		return fmt.Errorf("failed to save message in database: %w", err)
	}
	message, err := c.repo.Find(message)
	if err != nil {
		return fmt.Errorf("failed to find message: %w", err)
	}
	thread, err := c.repo.MessageThread(message)
	if err != nil {
		return fmt.Errorf("failed to fetch thread: %w", err)
	}
	chatCompletionThread := fromThread(thread)
	answer, err := c.service.Answer(cc.Ctx, chatCompletionThread, model)
	if err != nil {
		return err
	}
	replyMessage := cc.ReplyText(answer)
	telegramReplyMessage := telegram_persistence.FromTelegramMessage(&replyMessage)

	_, err = c.uRepo.GetOrCreate(telegramReplyMessage.FromUserID.Int64, telegram_persistence.FromExternalTelegramUser(replyMessage.From, replyMessage.Chat))
	if err != nil {
		log.Err(err).Send()
		return err
	}
	if err := c.repo.Create(telegramReplyMessage); err != nil {
		return err
	}
	return err
}

func (c *ChatGPTCommand) Command() string {
	return "chatgpt"
}

func (c *ChatGPTCommand) Description() string {
	return "Ask ChatGPT, anything!"
}

func (c *ChatGPTCommand) Match(message tbot.Message) bool {
	if message.ReplyToMessage == nil {
		log.Debug().Msg("Not a reply to a message")
		return false
	}
	reply := message.ReplyToMessage
	tMessage := telegram_persistence.FromTelegramMessage(reply)
	tMessage, err := c.repo.Find(tMessage)
	if err != nil {
		log.Err(err).Send()
		return false
	}
	thread, err := c.repo.MessageThread(tMessage)
	if err != nil {
		log.Logger.Err(err).Send()
		return false
	}
	firstMessage := thread[len(thread)-1]
	return strings.HasPrefix(firstMessage.Text, "/"+c.Command())
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
		m.ImGenModel = imGenFlux
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
