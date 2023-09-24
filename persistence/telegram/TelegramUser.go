package telegram

import (
	"database/sql"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"steplems-bot/persistence/spotify"
	"time"
)

type User struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	TelegramExternalID int64 `gorm:"primaryKey"`

	IsBot        bool   `json:"is_bot,omitempty"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name,omitempty"`
	UserName     string `json:"username,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`

	SpotifyUserID sql.NullString `gorm:"column:spotify_user_id"`
	SpotifyUser   spotify.User   `gorm:"foreignKey:SpotifyUserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (User) TableName() string {
	return "TelegramUser"
}

func FromExternalTelegramUser(user *tbot.User) User {
	return User{
		TelegramExternalID: user.ID,

		IsBot:        user.IsBot,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		UserName:     user.UserName,
		LanguageCode: user.LanguageCode,
	}
}
