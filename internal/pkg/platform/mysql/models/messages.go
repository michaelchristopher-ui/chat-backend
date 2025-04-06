package models

// Messages is the model for the Messsages Table
type Messages struct {
	ID         string `gorm:"id;primaryKey"`
	Message    string `gorm:"message"`
	ToUserID   string `gorm:"to_user_id"`
	Type       int    `gorm:"message_type"`
	FromUserID string `gorm:"from_user_id"`
	Timestamp  string `gorm:"timestamp"`
}
