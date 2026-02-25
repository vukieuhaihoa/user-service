package model

// User represents a user in the system.
// It maps to the "users" table in the database.
//
// Fields:
//   - ID: The unique identifier for the user (UUID).
//   - Username: The username of the user (unique, not null).
//   - Email: The email address of the user (unique, not null).
//   - Password: The hashed password of the user (not null).
//   - DisplayName: The display name of the user.
//   - CreatedAt: The timestamp when the user was created.
//   - UpdatedAt: The timestamp when the user was last updated.
type User struct {
	Base
	Username    string `gorm:"unique;not null;column:username" json:"username"`
	Email       string `gorm:"unique;not null;column:email" json:"email"`
	Password    string `gorm:"not null;column:password" json:"-"`
	DisplayName string `gorm:"column:display_name" json:"display_name"`
}

// TableName specifies the table name for the User model.
//
// Returns:
//   - string: The name of the database table for the User model
func (User) TableName() string {
	return "users"
}
