package chatservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"websocket_client/internal/pkg/core/adapter/chatadapter"
)

// getActualIP is a function that obtains IP the user is connected to from a databases
func (c ChatService) getActualIP(userID string) (string, error) {
	return c.Redis.GetValue(userID)
}

/*
Publish Message is a function that:
1. Checks whether the user IP is available, indicating that the user may be online
2. Attempts to send the message to the recipient through a request to a specific HTTP endpoint
3. Returns an error if any error is encountered during the entire process

An error is expected to not be returned if the user is not online when the request to the HTTP endpoint is made
*/
func (c ChatService) PublishMessage(req chatadapter.PublishMessageReq) error {
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
