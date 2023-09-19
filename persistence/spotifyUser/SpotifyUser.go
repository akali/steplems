package spotifyUser

import "gorm.io/gorm"

type SpotifyUser struct {
	gorm.Model
	Username       string
	Token          string
	RefresherToken string
}
