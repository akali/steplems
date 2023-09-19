package providers

import (
	"github.com/google/wire"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"steplems-bot/types"
)

func ProvideDatabaseConnectionURL() (types.DatabaseConnectionURL, error) {
	return ProvideEnvironmentVariable[types.DatabaseConnectionURL]("DATABASE")()
}

func ProvideDatabase(url types.DatabaseConnectionURL) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(string(url)), &gorm.Config{})
}

var DBProviders = wire.NewSet(ProvideDatabase, ProvideDatabaseConnectionURL)
