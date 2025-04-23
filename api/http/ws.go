package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"websocket_client/api/http/structs"
	"websocket_client/internal/pkg/core/adapter/chatadapter"
	"websocket_client/internal/pkg/core/adapter/wsconnadapter"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(c echo.Context) error {
	return c.JSON(http.StatusBadRequest, structs.ErrorRet{
		Error: "Something wrong occured",
	})
}

// Messages is a struct that defines the contents of the returned messages
type Message struct {
	FromUserID string `json:"from_user_id"`
	ToUserID   string `json:"to_user_id"`
	Message    string `json:"message"`
	Type       int    `json:"type"`
	Timestamp  string `json:"timestamp"`
}

// GetChatHistoryRes is a struct that defines the contents of the websocket request message sent to the user after executing the getChatHistory function
type GetChatHistoryReq struct {
	FromUserID     string `json:"from_user_id"`
	ToUserID       string `json:"to_user_id"`
	Offset         int    `json:"offset"`
	Limit          int    `json:"limit"`
	TimestampAfter string `json:"timestamp_after"`
}

// GetChatHistoryRes is a struct that defines the contents of the websocket response message sent to the user after executing the getChatHistory function
type GetChatHistoryRes struct {
	Messages []Message `json:"messages"`
	Error    string    `json:"error"`
}

// AddFriendReq is the json tagged request struct of the add friend feature
type AddFriendReq struct {
	FriendID string `json:"friend_id"`
}

// AddFriendResp is the json tagged response struct of the add friend feature
type AddFriendResp struct {
	Error string `json:"error"`
}
type SendMessageReq struct {
	ID         string `json:"id"`
	Message    string `json:"message"`
	Type       int    `json:"type"`
	ToUserID   string `json:"to_user_id"`
	FromUserID string `json:"from_user_id"`
	Timestamp  string `json:"timestamp"`
}

// RemoveFriendReq is the json tagged request struct of the remove friend feature
type RemoveFriendReq struct {
	FriendID string `json:"friend_id"`
}

type SearchUserReq struct {
	SearchUserID string `json:"search_user_id"`
}

/*
CallChatServiceFunc is a function that calls the chat service functions according to the supplied reqType,
supplying the other parameters as request parameters once it has been put into the correct request structs.
The function does return an error, but also sends a response to the user through the connection
*/
func (integrator *APIIntegrator) CallChatServiceFunc(dataByte []byte, userID, reqType string) error {
	var err error
	defer func() {
		if err != nil {
			integrator.Logger.NewError("err: %s", err.Error())
		}
	}()
	/*
		Nearly all the features that a user can request for is done through websocket.
		The switch handles it.
	*/
	var res any
	switch reqType {
	// Handles the request to send a message to another user
	case incomingMessageTypeMessage:
		req := SendMessageReq{}
		err = json.Unmarshal(dataByte, &req)
		if err != nil {
			return err
		}
		integrator.ChatService.SendMessage(chatadapter.SendMessageReq{
			ID:         req.ID,
			Message:    req.Message,
			Type:       req.Type,
			ToUserID:   req.ToUserID,
			FromUserID: req.FromUserID,
			Timestamp:  req.Timestamp,
		})

	// Handles the add friend request
	case incomingMessageTypeAddFriend:
		req := AddFriendReq{}
		err = json.Unmarshal(dataByte, &req)
		if err != nil {
			return err
		}

		integrator.ChatService.AddFriend(chatadapter.AddFriendReq{
			UserID:   userID,
			FriendID: req.FriendID,
		})

	// Handles the Get Chat History between the current user and another user
	case incomingMessageTypeGetChatHistory:
		req := GetChatHistoryReq{}
		err = json.Unmarshal(dataByte, &req)
		if err != nil {
			return err
		}

		res = integrator.ChatService.GetChatHistory(chatadapter.GetChatHistoryReq{
			FromUserID:     req.FromUserID,
			ToUserID:       req.ToUserID,
			Offset:         req.Offset,
			Limit:          req.Limit,
			TimestampAfter: req.TimestampAfter,
		})

	// Handles the remove friend request
	case incomingMessageTypeRemoveFriend:
		req := RemoveFriendReq{}
		err = json.Unmarshal(dataByte, &req)
		if err != nil {

			return err
		}
		res = integrator.ChatService.RemoveFriend(chatadapter.RemoveFriendReq{
			FriendID: req.FriendID,
		})

	// Handles the search user request
	case incomingMessageTypeSearchUser:
		req := SearchUserReq{}
		err = json.Unmarshal(dataByte, &req)
		if err != nil {
			return err
		}
		res = integrator.ChatService.SearchUser(chatadapter.SearchUserReq{
			SearchUserID: req.SearchUserID,
		})

	// Default case logs invalid requests
	default:
		integrator.Logger.NewError("Invalid Request, request: %s", string(dataByte))
		return errors.New("invalid request")
	}

	// Attempt to send a response to the requesting user
	isOnline, err := integrator.ChatService.ReceiveMessage(userID, res)
	if !isOnline {
		return fmt.Errorf("user %s goes offline before receiving response: %v", userID, res)
	}
	return err
}

/*
HandleWebsocket is the handler method for the /ws api endpoint
which does the upgrading of the connection before passing it to a processing function
*/

func (integrator *APIIntegrator) HandleWebsocket(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	// Attempt to upgrade the connection
	var conn wsconnadapter.Adapter
	conn, err := integrator.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		integrator.Logger.NewError("Error when upgrading for user %s, error: %s", userID, err.Error())
		return err
	}

	err = integrator.WsProcessor.ProcessWebsocketAfterUpgrade(conn, userID, integrator.CallChatServiceFunc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, structs.ErrorRet{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, nil)
}
