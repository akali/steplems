package deepinfra

import (
	"context"
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

func (s *DeepInfraService) GenerateImageFlux(ctx context.Context, prompt string) ([][]byte, error) {
	resp, err := s.client.CreateImage(ctx, deepinfra.WithPrompt(prompt), deepinfra.WithModel(deepinfra.FluxDev))
	if err != nil {
		return nil, err
	}
	return resp.ImageBytes(), nil
}

func (s *DeepInfraService) GenerateImage(ctx context.Context, prompt string) ([]string, error) {
	resp, err := s.client.CreateImage(ctx, deepinfra.WithPrompt(prompt), deepinfra.WithModel(deepinfra.StabilitySDXL))
	if err != nil {
		return nil, err
	}
	return resp.Urls(), err
}

func (s *DeepInfraService) VoiceToText(audioFilePath string) (string, error) {
	resp, err := s.client.VoiceToText(audioFilePath)
	if err != nil {
		return "", err
	}
	return resp.Text, nil
}
