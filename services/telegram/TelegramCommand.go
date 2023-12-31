package telegram

import (
	"context"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/wire"
	"steplems-bot/services/telegram/commands"
	"steplems-bot/types"
)

type CommandMap struct {
	commands map[string]TelegramCommand
}

func NewCommandMap(
	authorizeSpotifyCommand *commands.AuthorizeSpotifyCommand,
	helpCommand *commands.HelpCommand,
	nowPlayingCommand *commands.NowPlayingCommand,
	chatGPTCommand *commands.ChatGPTCommand,
) *CommandMap {
	cmdList := []TelegramCommand{
		helpCommand,
		authorizeSpotifyCommand,
		nowPlayingCommand,
		chatGPTCommand,
	}
	cm := CommandMap{commands: make(map[string]TelegramCommand)}
	for _, command := range cmdList {
		cm.commands[command.Command()] = command
	}
	return &cm
}

var CommandMapProvider = wire.NewSet(
	NewCommandMap,
	commands.CommandsProvider)

type TelegramCommand interface {
	Run(context.Context, types.Sender, tbot.Update) error
	Command() string
	Description() string
}

func ToAPITelegramCommand(cmd TelegramCommand) tbot.BotCommand {
	return tbot.BotCommand{
		Command:     cmd.Command(),
		Description: cmd.Description(),
	}
}
