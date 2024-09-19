package commands

import (
	"fmt"
	"steplems-bot/lib"
	"steplems-bot/services/chatgpt"
	"time"

	"github.com/rs/zerolog/log"
)

type ReasonCommand struct {
	reasoner *chatgpt.ReasonService
}

func NewReasonCommand(service *chatgpt.ReasonService) *ReasonCommand {
	return &ReasonCommand{reasoner: service}
}

func (r ReasonCommand) Run(cc *lib.ChatContext) error {
	respondText := "..."
	msg := cc.ReplyText(respondText)
	startTime := time.Now()
	respondText = "__Thinking...__"
	msg = cc.EditMessage(msg, respondText)
	log.Debug().Interface("msg", msg).Send()
	log.Debug().Msg("About to call Reason method")
	updatesChan := r.reasoner.Reason(cc.Ctx, cc.Text())
	log.Debug().Msg("Reason method called, waiting for updates")
	stopOnNext := false
	for {
		result, more := <-updatesChan
		if !more {
			break
		}
		log.Debug().Interface("result", result).Msg("Got result in reason command")
		if err := result.Err; err != nil {
			log.Err(err).Send()
			continue
		}
		respondText = fmt.Sprintf("%s\n%s\n\n**Thinking time:** __%q__\n", respondText, result.Data.Content, result.Duration)
		msg = cc.EditMessage(msg, respondText)
		if stopOnNext {
			break
		}
		if !result.Data.HasNext() {
			stopOnNext = true
		}
	}
	respondText = fmt.Sprintf("%s\n\n**Total thinking time:** %s", respondText, time.Since(startTime))
	msg = cc.EditMessage(msg, respondText)

	return nil
}

func (r ReasonCommand) Command() string {
	return "reason"
}

func (r ReasonCommand) Description() string {
	return "same as /chatgpt, but uses reasoning"
}
