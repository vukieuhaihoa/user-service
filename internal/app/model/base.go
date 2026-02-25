// package model defines the data models used in the application.
// It includes the User model representing users in the system.
package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base defines common fields for all database models.
type Base struct {
	ID        string    `gorm:"primaryKey;type:uuid;column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// BeforeCreate is a GORM hook that is triggered before a new User record is created in the database.
// It generates a new UUID for the user ID if it is not already set.
//
// Parameters:
//   - tx: The GORM database transaction
//
// Returns:
//   - error: An error if UUID generation fails, otherwise nil
func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}
