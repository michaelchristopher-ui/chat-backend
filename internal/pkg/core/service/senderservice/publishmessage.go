package senderservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"websocket_client/internal/pkg/core/adapter/senderadapter"

	"github.com/go-redis/redis/v8"
)

// getActualIP is a function that obtains IP the user is connected to from the redis database
func (c SenderService) getActualIP(userID string) (string, error) {
	return c.Redis.GetValue(userID)
}

/*
publishMessage is a function that:
1. Checks whether the user IP is available, indicating that the user may be online
2. Attempts to send the message to the recipient through a request to a specific HTTP endpoint
3. Returns an error if any error is encountered during the entire process

An error is expected to not be returned if the user is not online when the request to the HTTP endpoint is made
*/
func (c SenderService) PublishMessage(req senderadapter.SendMessageReq) (isOnline bool, err error) {
	//Checks whether the user IP is available
	ip, err := c.getActualIP(req.ToUserID)
	if ip == "" && err == nil {
		err = fmt.Errorf(logIPEmpty, req.ToUserID)
	}

	//After Checking, return false if user is not online indicated by the nonexistence of the user data in redis.
	if err != nil && err == redis.Nil {
		c.Logger.NewError(fmt.Sprintf("[ChatService][publishMessage] user was not online when attempting to send message, err: %s", err.Error()))
		return false, nil
	}
	if err != nil {
		c.Logger.NewError(fmt.Sprintf("[ChatService][publishMessage] error when obtaining recipient ip, err: %s", err.Error()))
		return false, err
	}

	// Sets up the message payload
	msg := MessagePayload{
		Message:    req.Message,
		FromUserID: req.FromUserID,
		ToUserID:   req.ToUserID,
		Type:       req.Type,
		Timestamp:  req.Timestamp,
	}
	payload, _ := json.Marshal(msg)

	bytePayload := bytes.NewBuffer(payload)

	// Attempts to send the message to the recipient through a request to a specific HTTP endpoint
	request, err := http.NewRequest("POST", "http://"+ip+"/receive", bytePayload)
	if err != nil {
		c.Logger.NewError(fmt.Sprintf("[ChatService][publishMessage] error when creating request instance, err: %s", err.Error()))
		return false, err
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		c.Logger.NewError(fmt.Sprintf("[ChatService][publishMessage] error when executing request, err: %s", err.Error()))
		return false, err
	}

	// Closes the body to prevent further writing
	defer func() {
		if response.Body.Close() != nil {
			c.Logger.NewError("[ChatService][publishMessage] Response failed to be closed")
		}
	}()

	// Decodes the response
	data := map[string]interface{}{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		c.Logger.NewError(fmt.Sprintf("[ChatService][publishMessage] error when decoding response body, err: %s", err.Error()))
		return false, err
	}

	if _, ok := data["error"]; ok {
		if errStr, ok2 := data["error"].(string); ok2 {
			err = errors.New(errStr)
		} else {
			err = ErrorExistsButNotConvertible
		}
	}

	if _, ok := data["is_online"]; ok {
		// By default, isOnline will be false should the conversion fail, but we log it anyway.
		isOnline = data["is_online"].(bool)
	}

	c.Logger.NewInfo(fmt.Sprintf("[ChatService][publishMessage] Done publishing message with ID %s", req.ID))
	return isOnline, err
}
