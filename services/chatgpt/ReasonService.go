package chatgpt

import (
	"context"
	"encoding/json"
	"fmt"
	"steplems-bot/types"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
)

const (
	SystemPrompt = `You are an expert AI assistant that explains your reasoning step by step. For each step, provide a title that describes what you're doing in that step, along with the content. Decide if you need another step or if you're ready to give the final answer. Respond in JSON format with 'title', 'content', and 'next_action' (either 'continue' or 'final_answer') keys. USE AS MANY REASONING STEPS AS POSSIBLE. AT LEAST 3. BE AWARE OF YOUR LIMITATIONS AS AN LLM AND WHAT YOU CAN AND CANNOT DO. IN YOUR REASONING, INCLUDE EXPLORATION OF ALTERNATIVE ANSWERS. CONSIDER YOU MAY BE WRONG, AND IF YOU ARE WRONG IN YOUR REASONING, WHERE IT WOULD BE. FULLY TEST ALL OTHER POSSIBILITIES. YOU CAN BE WRONG. WHEN YOU SAY YOU ARE RE-EXAMINING, ACTUALLY RE-EXAMINE, AND USE ANOTHER APPROACH TO DO SO. DO NOT JUST SAY YOU ARE RE-EXAMINING. USE AT LEAST 3 METHODS TO DERIVE THE ANSWER. RESPOND WITH ONLY ONE JSON OUTPUT. NEVER RESPOND WITH EMPTY JSON FIELDS. USE BEST PRACTICES.

Example of a valid JSON response:
` + "```" + `json
{
"title": "Identifying Key Information",
"content": "To begin solving this problem, we need to carefully examine the given information and identify the crucial elements that will guide our solution process. This involves...",
"next_action": "continue"
}` + "```"
	AssistantPrompt = "Thank you! I will now think step by step following my instructions, starting at the beginning after decomposing the problem."
	UserLastPrompt  = "Please provide the final answer based on your reasoning above."
)

const (
	MaxStepCount = 25
	MaxTokenSize = 300
)

var Model = types.ModelStorage{Model: "meta-llama/Meta-Llama-3.1-405B-Instruct", Backend: "deepinfra", ImGenModel: "stablediffusion"}

//var Model = types.ModelStorage{Model: "gpt-4o-mini", Backend: "openai", ImGenModel: "stablediffusion"}

type ReasonService struct {
	service *ChatGPTService
}

func NewReasonService(service *ChatGPTService) *ReasonService {
	return &ReasonService{service: service}
}

type StepData struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	NextAction string `json:"next_action"`
}

func (d StepData) json() string {
	b, _ := json.Marshal(d)
	return string(b)
}

func (d StepData) HasNext() bool {
	return d.NextAction != "final_answer"
}

type StepResult struct {
	Data     StepData
	Duration time.Duration
	Err      error
}

func removeHeaders(response string) string {
	endHeaderToken := "{"
	begin := strings.Index(response, endHeaderToken)
	if begin == -1 {
		return response
	}
	response = response[begin:]

	return response
}

func (r *ReasonService) makeCall(ctx context.Context, thread []openai.ChatCompletionMessage) (StepData, error) {
	return retry.DoWithData(func() (StepData, error) {
		answer, err := r.service.Answer(ctx, thread, Model, WithMaxTokenSize(MaxTokenSize))
		if err != nil {
			return StepData{}, err
		}
		answer, err = r.parseReasonJson(ctx, Model, answer)
		if err != nil {
			return StepData{}, err
		}
		var stepData StepData
		decoder := json.NewDecoder(strings.NewReader(answer))
		if err := decoder.Decode(&stepData); err != nil {
			return stepData, fmt.Errorf("failed to parse json for %q: %w", answer, err)
		}
		log.Debug().Interface("step_data", stepData).Msg("got step data")
		return stepData, nil
	}, retry.Attempts(3))
}

const TimeLimit = time.Minute * 5

func (r *ReasonService) reason(ctx context.Context, prompt string, resultChan chan<- StepResult) {
	defer close(resultChan)
	thread := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: SystemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
		{
			Role:    openai.ChatMessageRoleAssistant,
			Content: AssistantPrompt,
		},
	}
	allStartTime := time.Now()

	for steps := 0; steps <= MaxStepCount; steps++ {
		if time.Now().Sub(allStartTime) > TimeLimit {
			break
		}
		startTime := time.Now()
		stepData, err := r.makeCall(ctx, thread)
		endTime := time.Now()
		diff := endTime.Sub(startTime)
		stepResult := StepResult{
			Data:     stepData,
			Err:      err,
			Duration: diff,
		}
		resultChan <- stepResult
		thread = append(thread, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: stepData.json(),
		})
		if !stepData.HasNext() {
			break
		}
	}
	thread = append(thread, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: UserLastPrompt,
	})
	startTime := time.Now()
	stepData, err := r.makeCall(ctx, thread)
	endTime := time.Now()
	diff := endTime.Sub(startTime)
	stepResult := StepResult{
		Err:      err,
		Duration: diff,
		Data:     stepData,
	}
	resultChan <- stepResult
}

func (r *ReasonService) Reason(ctx context.Context, prompt string) <-chan StepResult {
	resultChan := make(chan StepResult, MaxStepCount)
	log.Info().Str("prompt", prompt).Msg("got reason request")
	go r.reason(ctx, prompt, resultChan)
	return resultChan
}
