package providers

import (
	"github.com/google/wire"
	"github.com/rs/zerolog"
	"steplems-bot/types"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func ProvideDatabaseConnectionURL() (types.DatabaseConnectionURL, error) {
	return ProvideEnvironmentVariable[types.DatabaseConnectionURL]("DATABASE")()
}

func ProvideDatabase(url types.DatabaseConnectionURL, logger zerolog.Logger) (*gorm.DB, error) {
	result, err := gorm.Open(sqlite.Open(string(url)), &gorm.Config{})
	if err != nil {
		logger.Warn().Err(err).Msg("failed to open database")
		return nil, nil
	}
	return result, err
}

var DBProviders = wire.NewSet(ProvideDatabase, ProvideDatabaseConnectionURL)
