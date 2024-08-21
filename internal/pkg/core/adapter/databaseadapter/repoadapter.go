package databaseadapter

import (
	"websocket_client/internal/pkg/platform/mysql/models"

	"gorm.io/gorm"
)

// Repoadapter defines an interface for a persistent database
type RepoAdapter interface {
	GetUserFriends(req GetUserFriendsReq) ([]models.UserFriends, error)
	DoCustomTransaction(fc func(tx *gorm.DB) error) error
	AddFriend(req AddFriendReq) error
	GetChatHistory(req GetChatHistoryReq) ([]models.Messages, error)
	GetAccount(req GetAccountReq) (models.Account, error)
	SetAccount(req SetAccountReq) error
	RemoveFriend(req RemoveFriendReq) error
}

type GetAccountReq struct {
	UserId string
}

type SetAccountReq struct {
	UserID   string
	Password string
}

type GetUserFriendsReq struct {
	UserID   string
	FriendID string
}

type AddFriendReq struct {
	UserID   string
	FriendID string
}

type RemoveFriendReq struct {
	UserID   string
	FriendID string
}

type StoreChatHistoryReq struct {
	Message    string
	ToUserID   string
	Type       int
	FromUserID string
	Timestamp  string
}

type GetChatHistoryReq struct {
	ToUserID       string
	FromUserID     string
	Offset         int
	Limit          int
	TimestampAfter string
}

type SendFunctionParam struct {
	Message    string `json:"message"`
	Type       int    `json:"type"`
	ToUserID   string `json:"to_user_id"`
	FromUserID string `json:"-"`
	Timestamp  string `json:"-"`
}
