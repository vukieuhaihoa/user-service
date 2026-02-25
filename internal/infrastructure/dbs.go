package infrastructure

import (
	"github.com/redis/go-redis/v9"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/common"
	redisPkg "github.com/vukieuhaihoa/bookmark-libs/pkg/redis"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/sqldb"

	"gorm.io/gorm"
)

// CreateRedisCon initializes and returns a Redis client.
func CreateRedisCon() *redis.Client {
	redisClient, err := redisPkg.NewClient("")
	common.HandlerError(err)

	return redisClient
}

// CreateSQLDBAndMigration initializes the SQL database client and performs migrations.
// It returns the initialized GORM DB instance.
// Returns:
//   - *gorm.DB: A pointer to the initialized GORM DB instance
func CreateSQLDBAndMigration() *gorm.DB {
	dbClient, err := sqldb.NewClient("")
	common.HandlerError(err)

	err = MigrateDB(dbClient)
	common.HandlerError(err)

	return dbClient
}

// CreateSQLDB initializes and returns a GORM DB client without performing migrations.
// It returns the initialized GORM DB instance.
// Returns:
//   - *gorm.DB: A pointer to the initialized GORM DB instance
func CreateSQLDB() *gorm.DB {
	dbClient, err := sqldb.NewClient("")
	common.HandlerError(err)

	return dbClient
}

// MigrateDB performs database migrations for the provided GORM DB instance.
// It migrates the necessary models to ensure the database schema is up to date.
//
// Parameters:
//   - db: A pointer to the GORM DB instance to perform migrations on
//
// Returns:
//   - error: An error object if the migration fails, otherwise nil
func MigrateDB(db *gorm.DB) error {
	return sqldb.MigrateSQLDB(db, "file://./migrations", "up", 0)
}
