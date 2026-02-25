package fixture

import (
	"testing"
	"time"

	"github.com/vukieuhaihoa/bookmark-libs/pkg/sqldb"
	"gorm.io/gorm"
)

var TestTime = time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)

// Fixture defines the interface for setting up test fixtures in the database.
type Fixture interface {
	// SetupDB sets up the database connection for the fixture.
	//
	// Parameters:
	//   - db: The GORM database connection to be used for the fixture
	SetupDB(db *gorm.DB)

	// Migrate migrates the database schema for the fixture.
	//
	// Returns:
	//   - error: An error if migration fails, otherwise nil
	Migrate() error

	// GenerateData generates test data for the fixture.
	//
	// Returns:
	//   - error: An error if data generation fails, otherwise nil
	GenerateData() error
	DB() *gorm.DB
}

// base provides common functionality for all test fixtures.
type base struct {
	db *gorm.DB
}

// SetupDB sets up the database connection for the base fixture.
//
// Parameters:
//   - db: The GORM database connection to be used for the fixture
func (b *base) SetupDB(db *gorm.DB) {
	b.db = db
}

// DB returns the GORM database connection used by the base fixture.
//
// Returns:
//   - *gorm.DB: The GORM database connection
func (b *base) DB() *gorm.DB {
	return b.db
}

// NewFixture initializes the test fixture by setting up the database,
// migrating the schema, and generating test data.
//
// Parameters:
//   - t: The testing object used for reporting errors
//   - fix: The fixture to be initialized
//
// Returns:
//   - *gorm.DB: A gorm.DB instance connected to the initialized test database
func NewFixture(t *testing.T, fix Fixture) *gorm.DB {
	// step 1: create test database
	fix.SetupDB(sqldb.InitMockDB(t))

	// step 2: migrate schema
	err := fix.Migrate()
	if err != nil {
		t.Fatalf("Failed to migrate db for testing")
	}
	// step 3: generate test data
	err = fix.GenerateData()
	if err != nil {
		t.Fatalf("Failed to generate test data: %v", err)
	}

	return fix.DB()
}
