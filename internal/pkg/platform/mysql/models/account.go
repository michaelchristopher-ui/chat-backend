package models

type Account struct {
	UserID   string `gorm:"userid;primaryKey"`
	Password string `gorm:"password;"`
}
