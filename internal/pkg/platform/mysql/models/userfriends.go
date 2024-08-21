package models

type UserFriends struct {
	UserID       string `gorm:"userid;primaryKey"`
	UserFriendID string `gorm:"userfriendid;primaryKey"`
}
