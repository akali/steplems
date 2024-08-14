package commands

import (
	"steplems-bot/lib"
	"steplems-bot/services/stablediffusion"
)

type StableDiffusionCommand struct {
	service *stablediffusion.StableDiffusionService
}

func NewStableDiffusionCommand(service *stablediffusion.StableDiffusionService) *StableDiffusionCommand {
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
	return "stablediffusion"
}

func (h *StableDiffusionCommand) Description() string {
	return "Generate image by prompt. Usage /stablediffusion An astronaut riding a rainbow unicorn."
}
