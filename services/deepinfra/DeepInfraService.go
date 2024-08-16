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

func (s *DeepInfraService) GenerateImageFlux(ctx context.Context, prompt string) ([]byte, error) {
	fluxResp, err := s.client.GenerateImageFlux(prompt)
	if err != nil {
		return nil, err
	}
	if len(fluxResp.Images) == 0 {
		return nil, fmt.Errorf("failed to generate image: no images returned")
	}
	return []byte(fluxResp.Images[0]), nil
}

func (s *DeepInfraService) GenerateImage(ctx context.Context, prompt string) (string, error) {
	resp, err := s.client.GenerateImageStableDiffusion(
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
