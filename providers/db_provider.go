package providers

import (
	"github.com/google/wire"
	"github.com/olehan/kek"
	"gorm.io/driver/postgres"

	"steplems-bot/types"

	"gorm.io/gorm"
)

func ProvideDatabaseConnectionURL() (types.DatabaseConnectionURL, error) {
	return ProvideEnvironmentVariable[types.DatabaseConnectionURL]("DATABASE")()
}

func ProvideDatabase(url types.DatabaseConnectionURL, lf *kek.Factory) (*gorm.DB, error) {
	logger := lf.NewLogger("DBProvider")
	result, err := gorm.Open(postgres.Open(string(url)), &gorm.Config{})
	if err != nil {
		logger.Warn.Println("failed to open database", err)
		return nil, nil
	}
	return result, err
}

var DBProviders = wire.NewSet(ProvideDatabase, ProvideDatabaseConnectionURL)
