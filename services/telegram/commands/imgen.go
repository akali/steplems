package commands

import (
	"fmt"
	"github.com/vincent-petithory/dataurl"
	"steplems-bot/lib"
	"steplems-bot/services/deepinfra"
)

type ImGenCommand struct {
	service *deepinfra.DeepInfraService
}

func NewImGenCommand(service *deepinfra.DeepInfraService) *ImGenCommand {
	return &ImGenCommand{
		service: service,
	}
}

func (h *ImGenCommand) Run(cc *lib.ChatContext) error {
	prompt := cc.Text()
	if model.ImGenModel == "flux" {
		b, err := h.service.GenerateImageFlux(cc.Ctx, prompt)
		if err != nil {
			return err
		}
		dataURL, err := dataurl.DecodeString(string(b))
		if err != nil {
			return err
		}
		cc.RespondImage(dataURL.Data)
		return nil
	}
	if model.ImGenModel == "stablediffusion" {
		url, err := h.service.GenerateImage(cc.Ctx, prompt)
		if err != nil {
			return err
		}
		cc.RespondImageURL(url)
		return nil
	}

	return nil
}

func (h *ImGenCommand) Command() string {
	return "imgen"
}

func (h *ImGenCommand) Description() string {
	return fmt.Sprintf("Generate image by prompt. Usage /%s An astronaut riding a rainbow unicorn.", h.Command())
}
