package spotify

import (
	"steplems-bot/persistence/spotifyUser"
)

type SpotifyService struct {
	userRepo *spotifyUser.SpotifyUserRepository
}

func NewSpotifyService(userRepo *spotifyUser.SpotifyUserRepository) *SpotifyService {
	return &SpotifyService{userRepo: userRepo}
}

func (s *SpotifyService) AuthorizeUser(username string) (spotifyUser.SpotifyUser, error) {
	user, err := s.userRepo.Create(spotifyUser.SpotifyUser{Username: username})
	if err != nil {
		return spotifyUser.SpotifyUser{}, err
	}
	return user, nil
}

func (s *SpotifyService) FindAll() []spotifyUser.SpotifyUser {
	return s.userRepo.FindAll()
}
