package telegram_persistence

import (
	"database/sql"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"time"
)

type Message struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	MessageID int   `gorm:"primaryKey;autoIncrement:false"`
	ChatID    int64 `gorm:"primaryKey;autoIncrement:false"`
	Text      string

	ReplyMessageID sql.NullInt32 `gorm:"column:reply_message_id"`
	ReplyChatID    sql.NullInt64 `gorm:"column:reply_chat_id"`
	ReplyMessage   *Message      `gorm:"foreignKey:ReplyMessageID,ReplyChatID;references:MessageID,ChatID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	FromUserID sql.NullInt64 `gorm:"column:from_user_id"`
	FromUser   *User         `gorm:"foreignKey:FromUserID;references:TelegramExternalID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (Message) TableName() string {
	return "TelegramMessage"
}

func FromTelegramMessage(message *tbot.Message) Message {
	msg := Message{
		MessageID:  message.MessageID,
		ChatID:     message.Chat.ID,
		Text:       message.Text,
		FromUserID: sql.NullInt64{Int64: message.From.ID, Valid: true},
	}
	if message.ReplyToMessage != nil {
		msg.ReplyMessageID = sql.NullInt32{Int32: int32(message.ReplyToMessage.MessageID), Valid: true}
		msg.ReplyChatID = sql.NullInt64{Int64: message.ReplyToMessage.Chat.ID, Valid: true}
	}
	return msg
}
