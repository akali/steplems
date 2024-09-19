package telegram

import (
	tbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/wire"
	"steplems-bot/lib"
	"steplems-bot/services/telegram/commands"
	"strings"
)

type CommandMap struct {
	commands       map[string]TelegramCommand
	chatGPTCommand *commands.ChatGPTCommand
}

func (m CommandMap) Get(message tbot.Message) (TelegramCommand, bool) {
	if message.IsCommand() {
		if command, ok := m.commands[message.Command()]; ok {
			return command, ok
		}
		return nil, false
	}
	if m.chatGPTCommand.Match(message) {
		return m.chatGPTCommand, true
	}
	return nil, false
}

func (m CommandMap) Match(message tbot.Message) bool {
	match := m.chatGPTCommand.Match(message)
	return match
}

func (m CommandMap) ApiCommands() []tbot.BotCommand {
	var cmds []tbot.BotCommand
	for _, command := range m.commands {
		cmd := ToAPITelegramCommand(command)
		cmd.Command = strings.TrimPrefix(strings.Split(cmd.Command, " ")[0], "/")
		cmds = append(cmds, cmd)
	}
	return cmds
}

func NewCommandMap(
	authorizeSpotifyCommand *commands.AuthorizeSpotifyCommand,
	helpCommand *commands.HelpCommand,
	nowPlayingCommand *commands.NowPlayingCommand,
	chatGPTCommand *commands.ChatGPTCommand,
	setModelCommand *commands.SetModelCommand,
	stableDiffusionCommand *commands.ImGenCommand,
	transcribeCommand *commands.TranscribeCommand,
	reasonCommand *commands.ReasonCommand,
) *CommandMap {
	cmdList := []TelegramCommand{
		helpCommand,
		authorizeSpotifyCommand,
		nowPlayingCommand,
		chatGPTCommand,
		setModelCommand,
		stableDiffusionCommand,
		transcribeCommand,
		reasonCommand,
	}
	cm := CommandMap{commands: make(map[string]TelegramCommand)}
	for _, command := range cmdList {
		cm.commands[command.Command()] = command
	}
	cm.chatGPTCommand = chatGPTCommand
	return &cm
}

var CommandMapProvider = wire.NewSet(
	NewCommandMap,
	commands.CommandsProvider)

type TelegramCommand interface {
	Run(cc *lib.ChatContext) error
	Command() string
	Description() string
}

func ToAPITelegramCommand(cmd TelegramCommand) tbot.BotCommand {
	return tbot.BotCommand{
		Command:     cmd.Command(),
		Description: cmd.Description(),
	}
}
