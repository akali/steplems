package deepinfra

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"steplems-bot/lib/deepinfra"
)

type DeepInfraService struct {
	client *deepinfra.Client
	logger zerolog.Logger
}

func NewStableDiffusionService(client *deepinfra.Client, logger zerolog.Logger) *DeepInfraService {
	return &DeepInfraService{client: client, logger: logger}
}

func (s *DeepInfraService) GenerateImage(ctx context.Context, prompt string) (string, error) {
	resp, err := s.client.GenerateImage(
		&deepinfra.Request{
			Input: &deepinfra.Input{
				Prompt: &prompt,
			},
		},
	)
	if err != nil {
		return "", err
	}
	if resp.Output == nil || len(resp.Output) == 0 {
		return "", fmt.Errorf("got empty response")
	}
	return resp.Output[0], nil
}

func (s *DeepInfraService) VoiceToText(audioFilePath string) (string, error) {
	resp, err := s.client.VoiceToText(audioFilePath)
	if err != nil {
		return "", err
	}
	return resp.Text, nil
}
