package mysql

import (
	"websocket_client/internal/conf"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
	"websocket_client/internal/pkg/platform/mysql/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

func NewDatabase() (databaseadapter.RepoAdapter, error) {
	dsn := conf.GetConfig().Database.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.Set("gorm:websocket_server", "ENGINE=InnoDB").AutoMigrate(
		models.ModelsList...,
	)
	if err != nil {
		return nil, err
	}
	return Database{
		db: db,
	}, nil
}

func (d Database) DoCustomTransaction(fc func(tx *gorm.DB) error) error {
	return d.db.Transaction(fc)
}

// GetUserFriends obtains users that are friends with the supplied UserID
func (d Database) GetUserFriends(req databaseadapter.GetUserFriendsReq) (users []models.UserFriends, err error) {
	err = d.db.Where(models.UserFriends{
		UserID: req.UserID,
	}).Find(&users).Error
	return
}

// AddFriend adds two entries so the relationship user1 <-> user2 is established within the DB.
func (d Database) AddFriend(req databaseadapter.AddFriendReq) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		userfriends := models.UserFriends{
			UserID:       req.UserID,
			UserFriendID: req.FriendID,
		}
		if err := tx.Create(&userfriends).Error; err != nil {
			return err
		}
		userfriends = models.UserFriends{
			UserID:       req.FriendID,
			UserFriendID: req.UserID,
		}
		return tx.Create(&userfriends).Error
	})
}

// RemoveFriend removes two entries so the relationship user1 <-> user2 gets removed within the DB.
func (d Database) RemoveFriend(req databaseadapter.RemoveFriendReq) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		userfriends := models.UserFriends{
			UserID:       req.UserID,
			UserFriendID: req.FriendID,
		}
		if err := tx.Delete(&userfriends).Error; err != nil {
			return err
		}
		userfriends = models.UserFriends{
			UserID:       req.FriendID,
			UserFriendID: req.UserID,
		}
		return tx.Delete(&userfriends).Error
	})
}

// GetChatHistory obtains messages between two users with limit and offset for pagination purposes.
func (d Database) GetChatHistory(req databaseadapter.GetChatHistoryReq) (messages []models.Messages, err error) {
	err = d.db.Where(models.Messages{
		ToUserID:   req.ToUserID,
		FromUserID: req.FromUserID,
	}).Or(d.db.Where(models.Messages{
		ToUserID:   req.FromUserID,
		FromUserID: req.ToUserID,
	})).Where("timestamp > ?",
		req.TimestampAfter).Order("timestamp desc").Limit(req.Limit).Offset(req.Offset).Find(&messages).Error

	return
}

func (d Database) GetAccount(req databaseadapter.GetAccountReq) (account models.Account, err error) {
	err = d.db.Where(models.Account{
		UserID: req.UserId,
	}).Find(&account).Error
	return
}

func (d Database) SetAccount(req databaseadapter.SetAccountReq) error {
	account := models.Account{
		UserID:   req.UserID,
		Password: req.Password,
	}
	return d.db.Create(&account).Error
}
