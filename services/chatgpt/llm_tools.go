package chatgpt

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"os"
	"steplems-bot/types"
)

func mustReadPrompt(name string) string {
	filename := fmt.Sprintf("static/prompts/%s", name)
	b, err := os.ReadFile(filename)
	if err != nil {
		panic(fmt.Errorf("file %s not found", filename))
	}
	return string(b)
}

func (r *ReasonService) parseReasonJson(ctx context.Context, model types.ModelStorage, json string) (string, error) {
	thread := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: mustReadPrompt("json_system_prompt.txt"),
		}, {
			Role:    openai.ChatMessageRoleUser,
			Content: mustReadPrompt("json_user_prompt_1.txt"),
		}, {
			Role:    openai.ChatMessageRoleAssistant,
			Content: mustReadPrompt("json_assistant_prompt_1.txt"),
		}, {
			Role:    openai.ChatMessageRoleUser,
			Content: mustReadPrompt("json_user_prompt_2.txt"),
		}, {
			Role:    openai.ChatMessageRoleAssistant,
			Content: mustReadPrompt("json_assistant_prompt_2.txt"),
		}, {
			Role:    openai.ChatMessageRoleUser,
			Content: mustReadPrompt("json_user_prompt_3.txt"),
		}, {
			Role:    openai.ChatMessageRoleAssistant,
			Content: mustReadPrompt("json_assistant_prompt_3.txt"),
		}, {
			Role:    openai.ChatMessageRoleUser,
			Content: json,
		},
	}
	return r.service.Answer(ctx, thread, model)
}
