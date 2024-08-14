package commands

import (
	"steplems-bot/lib"
)

type HelpCommand struct{}

func (h *HelpCommand) Run(cc *lib.ChatContext) error {
	cc.RespondText("help requested")
	return nil
}

func (h *HelpCommand) Command() string {
	return "help"
}

func (h *HelpCommand) Description() string {
	return "Get help."
}

func NewHelpCommand() *HelpCommand {
	return &HelpCommand{}
}
