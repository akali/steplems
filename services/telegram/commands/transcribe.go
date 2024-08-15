package commands

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"steplems-bot/lib"
	"steplems-bot/services/deepinfra"
)

type TranscribeCommand struct {
	service *deepinfra.DeepInfraService
}

func (t *TranscribeCommand) Run(cc *lib.ChatContext) error {
	reply := cc.Update.Message.ReplyToMessage
	if reply == nil {
		cc.RespondText(fmt.Sprintf("no reply, usage: reply with /%s command on voiceMessage message", t.Command()))
		return nil
	}
	voiceMessage := reply.Voice
	if voiceMessage == nil {
		cc.RespondText(fmt.Sprintf("no voiceMessage on reply, usage: reply with /%s command on voiceMessage message", t.Command()))
		return nil
	}

	bVoice, err := cc.GetFile(voiceMessage.FileID)
	if err != nil {
		return err
	}
	folder, err := os.MkdirTemp("/tmp", "yt*")
	if err != nil {
		return err
	}
	filename := path.Join(folder, voiceMessage.FileUniqueID)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, bytes.NewBuffer(bVoice))
	if err != nil {
		return err
	}
	text, err := t.service.VoiceToText(filename)
	if err != nil {
		return err
	}
	cc.RespondText(text)
	return nil
}

func (t *TranscribeCommand) Command() string {
	return "transcribe"
}

func (t *TranscribeCommand) Description() string {
	return "transcribes given audio"
}

func NewTranscribeCommand(service *deepinfra.DeepInfraService) *TranscribeCommand {
	return &TranscribeCommand{service: service}
}
