package telegram_persistence

import (
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"steplems-bot/persistence/spotify_persistence"
)

type UserRepository struct {
	DB *gorm.DB
}

func (p *UserRepository) RunMigrations() error {
	return p.DB.AutoMigrate(&User{})
}

func NewUserRepository(DB *gorm.DB) *UserRepository {
	return &UserRepository{DB}
}

func (p *UserRepository) FindAll() []User {
	var result []User
	p.DB.Find(&result)
	return result
}

func (p *UserRepository) Create(user User) (User, error) {
	result := p.DB.Create(&user)
	return user, result.Error
}

func (p *UserRepository) Get(externalID int64) (User, error) {
	var user User
	result := p.DB.Where("telegram_external_id = ?", externalID).First(&user)
	if result.Error != nil {
		return User{}, result.Error // User not found or other error
	}
	return user, nil // User found and returned
}

func (p *UserRepository) GetOrCreate(externalID int64, newUser User) (User, error) {
	// Try to find the user by external ID
	var existingUser User
	result := p.DB.Where("telegram_external_id = ?", externalID).First(&existingUser)

	if result.Error == nil {
		// User already exists, return the existing user
		return existingUser, nil
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// User doesn't exist, so create a new one
		return newUser, p.DB.Create(&newUser).Error
	}
	return User{}, result.Error // Other error occurred
}

var NoSpotifyUserFound = errors.New("no spotify user found")

func (p *UserRepository) EnsureSpotifyUserExists(telegramUserID int64) (spotify_persistence.User, error) {
	// Check if the TelegramUser already has a SpotifyUser
	var telegramUser User
	result := p.DB.Preload("SpotifyUser").Where("telegram_external_id = ?", telegramUserID).First(&telegramUser)
	if result.Error != nil {
		return spotify_persistence.User{}, result.Error // Error occurred while checking
	}

	// If the TelegramUser does not have a SpotifyUser, return an error
	if telegramUser.SpotifyUser.ID == "" {
		return spotify_persistence.User{}, NoSpotifyUserFound
	}

	// Return the existing SpotifyUser
	return telegramUser.SpotifyUser, nil
}

func (p *UserRepository) SaveSpotifyUser(user User, spotifyUser spotify_persistence.User) error {
	user.SpotifyUserID = sql.NullString{String: spotifyUser.ID, Valid: true}
	return p.DB.Save(&user).Error
}
