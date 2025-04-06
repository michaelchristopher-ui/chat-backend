package models

// Account is the model for the Account Table
type Account struct {
	UserID   string `gorm:"userid;primaryKey"`
	Password string `gorm:"password;"`
}
