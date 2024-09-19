package telegram_persistence

import (
	"fmt"
	"gorm.io/gorm"
)

type MessageRepository struct {
	DB *gorm.DB
}

func (p *MessageRepository) RunMigrations() error {
	return p.DB.AutoMigrate(&Message{})
}

func NewMessageRepository(DB *gorm.DB) *MessageRepository {
	return &MessageRepository{
		DB: DB,
	}
}

func (p *MessageRepository) MessageThread(message Message) ([]Message, error) {
	var result []Message

	for message.ReplyMessageID.Valid {
		result = append(result, message)
		var err error
		if message, err = p.Find(Message{
			ChatID:    message.ReplyChatID.Int64,
			MessageID: int(message.ReplyMessageID.Int32),
		}); err != nil {
			return nil, fmt.Errorf("unexpected state, replying message not found in thread: %w", err)
		}
	}

	result = append(result, message)

	return result, nil
}

func (p *MessageRepository) Create(message Message) error {
	return p.DB.Create(&message).Error
}

func (p *MessageRepository) Find(message Message) (Message, error) {
	result := p.DB.Preload("FromUser").Where("message_id = ? AND chat_id = ?", message.MessageID, message.ChatID).First(&message)
	return message, result.Error
}
