package healthcheck

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/sqldb"
	"gorm.io/gorm"
)

func TestRepository_DBPing(t *testing.T) {
	t.Parallel()

	// Since DBPing is not implemented, we will just test that it returns nil for now.
	testCases := []struct {
		name string

		setupMockDB func() *gorm.DB

		expectedErrStr string
		expectedError  error
	}{
		{
			name: "DBPing returns nil",

			setupMockDB: func() *gorm.DB {
				db := sqldb.InitMockDB(t)

				return db
			},

			expectedError: nil,
		},
		{
			name: "DBPing on closed DB",

			setupMockDB: func() *gorm.DB {
				db := sqldb.InitMockDB(t)
				sqlDB, _ := db.DB()
				sqlDB.Close()

				return db
			},

			expectedError: sql.ErrConnDone,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			db := tc.setupMockDB()

			healthCheckRepo := NewHealthCheckRepository(nil, db)

			err := healthCheckRepo.DBPing(ctx)

			if tc.expectedError != nil {
				assert.Error(t, err)
				if tc.expectedErrStr != "" {
					assert.Contains(t, err.Error(), tc.expectedErrStr)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
