package chatservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"websocket_client/internal/pkg/core/adapter/databaseadapter"
)

/*
sendWebsocket is a function that:
1. Attempts to send a message to the user's websocket if the connection is available
2. Keep retrying when it fails to send the message as long as the user is connected
*/
func (c ChatService) sendWebsocket(userID string, data interface{}) {
	if conn, ok := c.UserConnections[userID]; ok {
		var err error
		for ok {
			err = conn.WriteJSON(data)
			if err == nil {
				break
			}
			conn, ok = c.UserConnections[userID]
		}
	}
}

// broadcastEmptyMessageToFriends is a function intended to send an empty message to the user's friends indicating the user's online/offline status
func (c ChatService) broadcastEmptyMessageToFriends(userID string, messageType int) {
	userFriends, err := c.DB.GetUserFriends(databaseadapter.GetUserFriendsReq{
		UserID: userID,
	})

	if err != nil {
		for _, eachUserFriend := range userFriends {
			publishMessageReq := PublishMessageReq{
				Message:    "",
				ToUserID:   eachUserFriend.UserFriendID,
				FromUserID: userID,
				Type:       messageType,
			}
			err := c.publishMessage(publishMessageReq)
			if err != nil {
				log.Printf(logErrMessageCannotPublish, err.Error(), publishMessageReq)
			}
		}
	}
}

// getActualIP is a function that obtains IP the user is connected to from a databases
func (c ChatService) getActualIP(userID string) (string, error) {
	return c.Redis.GetValue(userID)
}

type PublishMessageReq struct {
	Message    string `json:"message"`
	Type       int    `json:"type"`
	ToUserID   string `json:"to_user_id"`
	FromUserID string `json:"-"`
	Timestamp  string `json:"-"`
}

/*
publishMessage is a function that:
1. Checks whether the user IP is available, indicating that the user may be online
2. Attempts to send the message to the recipient through a request to a specific HTTP endpoint
3. Returns an error if any error is encountered during the entire process

An error is expected to not be returned if the user is not online when the request to the HTTP endpoint is made
*/
func (c ChatService) publishMessage(req PublishMessageReq) error {
	//Checks whether the user IP is available
	ip, err := c.getActualIP(req.ToUserID)
	if err != nil {
		return err
	}
	if ip == "" {
		return fmt.Errorf(logIPEmpty, ip)
	}

	//Attempts to send the message to the recipient through a request to a specific HTTP endpoint
	msg := MessagePayload{
		Message:    req.Message,
		FromUserID: req.FromUserID,
		Type:       req.Type,
		Timestamp:  req.Timestamp,
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	bytePayload := bytes.NewBuffer(payload)

	request, err := http.NewRequest("POST", "http://"+ip+"/receive", bytePayload)
	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	data := map[string]interface{}{}

	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return err
	}

	if _, ok := data["error"]; ok {
		if errStr, ok2 := data["error"].(string); ok2 {
			err = errors.New(errStr)
		} else {
			err = errorExistsButNotConvertible
		}
	}

	return err
}
