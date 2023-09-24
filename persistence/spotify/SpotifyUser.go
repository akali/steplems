package spotify

import (
	sapi "github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"time"
)

type User struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	ID           string    `gorm:"primaryKey"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
	TokenType    string    `json:"token_type,omitempty"`
}

func (User) TableName() string {
	return "SpotifyUser"
}

func (u User) OAuthToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  u.AccessToken,
		TokenType:    u.TokenType,
		RefreshToken: u.RefreshToken,
		Expiry:       u.Expiry,
	}
}

func PrivateUserToUser(user *sapi.PrivateUser) User {
	return User{
		ID: user.ID,
	}
}

func (u User) SetOAuthToken(token *oauth2.Token) User {
	u.AccessToken = token.AccessToken
	u.RefreshToken = token.RefreshToken
	u.Expiry = token.Expiry
	u.TokenType = token.TokenType
	return u
}
