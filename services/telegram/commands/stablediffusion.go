package commands

import (
	"steplems-bot/lib"
	"steplems-bot/services/deepinfra"
)

type StableDiffusionCommand struct {
	service *deepinfra.DeepInfraService
}

func NewStableDiffusionCommand(service *deepinfra.DeepInfraService) *StableDiffusionCommand {
	return &StableDiffusionCommand{
		service: service,
	}
}

func (h *StableDiffusionCommand) Run(cc *lib.ChatContext) error {
	prompt := cc.Text()
	url, err := h.service.GenerateImage(cc.Ctx, prompt)
	if err != nil {
		return err
	}
	cc.RespondImageURL(url)
	return nil
}

func (h *StableDiffusionCommand) Command() string {
	return "deepinfra"
}

func (h *StableDiffusionCommand) Description() string {
	return "Generate image by prompt. Usage /deepinfra An astronaut riding a rainbow unicorn."
}
