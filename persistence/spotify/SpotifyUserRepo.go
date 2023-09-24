package spotify

import (
	"gorm.io/gorm"
	"steplems-bot/persistence"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewSpotifyUserRepository(DB *gorm.DB) *UserRepository {
	return &UserRepository{DB}
}

func (p *UserRepository) FindAll() []User {
	return persistence.FindAll[User](p.DB)
}

func (p *UserRepository) Create(user User) (User, error) {
	result := p.DB.Create(&user)
	return user, result.Error
}

func (p *UserRepository) RunMigrations() error {
	err := p.DB.AutoMigrate(&User{})
	if err != nil {
		return err
	}
	constraint := "fk_TelegramUser_spotify_user"
	if p.DB.Migrator().HasConstraint(&User{}, constraint) {
		return p.DB.Migrator().DropConstraint(&User{}, constraint)
	}
	return nil
}
