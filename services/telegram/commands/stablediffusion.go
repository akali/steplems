package commands

import (
	"fmt"
	"steplems-bot/lib"
	"steplems-bot/services/deepinfra"
)

type ImGenCommand struct {
	service *deepinfra.DeepInfraService
}

func NewStableDiffusionCommand(service *deepinfra.DeepInfraService) *ImGenCommand {
	return &ImGenCommand{
		service: service,
	}
}

func (h *ImGenCommand) Run(cc *lib.ChatContext) error {
	prompt := cc.Text()
	url, err := h.service.GenerateImage(cc.Ctx, prompt, model.ImGenModel)
	if err != nil {
		return err
	}
	cc.RespondImageURL(url)
	return nil
}

func (h *ImGenCommand) Command() string {
	return "imgen"
}

func (h *ImGenCommand) Description() string {
	return fmt.Sprintf("Generate image by prompt. Usage /%s An astronaut riding a rainbow unicorn.", h.Command())
}
