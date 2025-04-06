package models

// UserFriends is a model for the UserFriends table
type UserFriends struct {
	UserID       string `gorm:"userid;primaryKey"`
	UserFriendID string `gorm:"userfriendid;primaryKey"`
}
