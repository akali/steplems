package types

type DatabaseConnectionURL string

type MigrationRunner interface {
	RunMigrations() error
}
