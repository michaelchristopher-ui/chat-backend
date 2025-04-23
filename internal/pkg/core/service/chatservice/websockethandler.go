package chatservice

import (
	"websocket_client/internal/common"
)

// SetValueOnRedis setup user and connected server key-value entry to redis so that services can find out which server the user is connected to
func (c ChatService) SetValueOnRedis(userID string) *bool {
	isOpen := true
	c.Redis.SetValueUntilChannelClose(userID, *common.IPPort, 30, &isOpen)
	return &isOpen
}
