package spotify_user

import (
	"gorm.io/gorm"
	"steplems-bot/persistence"
)

type SpotifyUserRepository struct {
	DB *gorm.DB
}

func NewSpotifyUserRepository(DB *gorm.DB) SpotifyUserRepository {
	return SpotifyUserRepository{DB}
}

func (p *SpotifyUserRepository) FindAll() []SpotifyUser {
	return persistence.FindAll[SpotifyUser](p.DB)
}

func (p *SpotifyUserRepository) Create(user SpotifyUser) (SpotifyUser, error) {
	result := p.DB.Create(&user)
	return user, result.Error
}

func (p *SpotifyUserRepository) FindByUsername(username string) *SpotifyUser {
	var spotifyUser SpotifyUser
	p.DB.Where(&SpotifyUser{Username: username}).First(&spotifyUser)
	return &spotifyUser
}
