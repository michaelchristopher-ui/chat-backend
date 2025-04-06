package databaseadapter

import (
	"websocket_client/internal/pkg/platform/mysql/models"

	"gorm.io/gorm"
)

// Repoadapter defines an interface for a persistent database
//
//go:generate mockgen -source=repoadapter.go -package=databaseadapter -destination=repoadapter_mock.go
type RepoAdapter interface {
	GetUserFriends(req GetUserFriendsReq) ([]models.UserFriends, error)
	DoCustomTransaction(fc func(tx *gorm.DB) error) error
	AddFriend(req AddFriendReq) error
	GetChatHistory(req GetChatHistoryReq) ([]models.Messages, error)
	GetAccount(req GetAccountReq) (models.Account, error)
	SetAccount(req SetAccountReq) error
	SaveChatHistoryWithTx(tx *gorm.DB, req SaveChatHistoryWithTxReq) error
	RemoveFriend(req RemoveFriendReq) error
	SearchFriend(req SearchFriendRequest) (users []models.Account, err error)
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

type SearchFriendRequest struct {
	UserID string
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

type SaveChatHistoryWithTxReq struct {
	ID         string
	Message    string
	ToUserID   string
	Type       int
	FromUserID string
	Timestamp  string
}
