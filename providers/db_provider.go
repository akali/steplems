package providers

import (
	"github.com/google/wire"
	"gorm.io/driver/postgres"

	"steplems-bot/types"

	"gorm.io/gorm"
)

func ProvideDatabaseConnectionURL() (types.DatabaseConnectionURL, error) {
	return ProvideEnvironmentVariable[types.DatabaseConnectionURL]("DATABASE")()
}

func ProvideDatabase(url types.DatabaseConnectionURL) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(string(url)), &gorm.Config{})
}

var DBProviders = wire.NewSet(ProvideDatabase, ProvideDatabaseConnectionURL)
