package mysql

import (
	"errors"
	"fmt"
	"websocket_client/internal/conf"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
	"websocket_client/internal/pkg/platform/mysql/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

// NewDatabase opens a new connection to a mySQL database using the GORM ORM.
func NewDatabase() (databaseadapter.RepoAdapter, error) {
	// Obtain dsn from the config
	dsn := conf.GetConfig().Database.Dsn
	// Attempt to open a connection according to the dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Attempt to automigrate the database according to the list of models used
	err = db.Set("gorm:websocket_server", "ENGINE=InnoDB").AutoMigrate(
		models.ModelsList...,
	)
	if err != nil {
		return nil, err
	}

	// Return the connection within the database struct
	return Database{
		db: db,
	}, nil
}

// DoCustomTransaction executes the Gorm Transaction
func (d Database) DoCustomTransaction(fc func(tx *gorm.DB) error) error {
	return d.db.Transaction(fc)
}

// GetUserFriends obtains users that are friends with the supplied UserID
func (d Database) GetUserFriends(req databaseadapter.GetUserFriendsReq) (userFriends []models.UserFriends, err error) {
	err = d.db.Where(models.UserFriends{
		UserID: req.UserID,
	}).Find(&userFriends).Error
	return
}

/*
SearchUser searches a user, using the UserID in the search as a substring
TODO: Move this to elasticsearch if we require a faster search
*/
func (d Database) SearchUser(req databaseadapter.SearchUserRequest) (users []models.Account, err error) {
	/*
		We add a percent sign as a prefix and suffix of the UserID for a substring search.
		Double percent sign escapes the percent sign in the string format
	*/
	wildcardedUserID := fmt.Sprintf("%%%s%%", req.UserID)
	/*
		Query the database with a LIKE query, supplying the wildcarded user ID as a parameter.
		We use string conditions since GORM does not have an interface for LIKE queries.
		The query is safe since GORM escapes it.
	*/
	err = d.db.Where("userid LIKE ?", wildcardedUserID).Find(&users).Error
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

// SaveChatHistory saves a message to the database with TX
func (d Database) SaveChatHistoryWithTx(tx *gorm.DB, req databaseadapter.SaveChatHistoryWithTxReq) (err error) {
	tx = tx.Create(&models.Messages{
		ID:         req.ID,
		ToUserID:   req.ToUserID,
		FromUserID: req.FromUserID,
		Type:       req.Type,
		Timestamp:  req.Timestamp,
	})

	// Gorm returns no error if rows are affected. We consider this to be an error.
	if tx.RowsAffected == 0 {
		return errors.New("no rows were updated during the saving of chat history")
	}

	return tx.Error
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

// GetAccount obtains account data from user id. Password is hashed with bcrypt.
func (d Database) GetAccount(req databaseadapter.GetAccountReq) (account models.Account, err error) {
	err = d.db.Where(models.Account{
		UserID: req.UserId,
	}).Find(&account).Error
	return
}

// SetAccount inserts a new user account entry to the database
func (d Database) SetAccount(req databaseadapter.SetAccountReq) error {
	account := models.Account{
		UserID:   req.UserID,
		Password: req.Password,
	}
	return d.db.Create(&account).Error
}
