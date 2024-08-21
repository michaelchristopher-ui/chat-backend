package models

type Messages struct {
	ID         int    `gorm:"id;primaryKey"`
	Message    string `gorm:"message"`
	ToUserID   string `gorm:"to_user_id"`
	Type       int    `gorm:"message_type"`
	FromUserID string `gorm:"from_user_id"`
	Timestamp  string `gorm:"timestamp"`
}
