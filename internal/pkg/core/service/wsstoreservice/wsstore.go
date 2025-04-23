package wsstoreservice

import (
	"fmt"
	"sync"
	"websocket_client/internal/pkg/core/adapter/wsconnadapter"
	"websocket_client/internal/pkg/core/adapter/wsstoreadapter"
)

type NewWsStoreServiceReq struct {
	UserConnections map[string]wsconnadapter.Adapter
	Lock            map[string]*sync.Mutex
}

type WsStore struct {
	UserConnections map[string]wsconnadapter.Adapter
	Lock            map[string]*sync.Mutex
	LockMutex       sync.Mutex
}

func NewChatBackendService(req NewWsStoreServiceReq) wsstoreadapter.Adapter {
	return &WsStore{
		UserConnections: req.UserConnections,
		Lock:            req.Lock,
	}
}

// AddConn adds a user connection key value entry in the online users map provided that it is able to obtain all required locks
func (c *WsStore) AddConn(conn wsconnadapter.Adapter, userID string) error {
	c.LockMutex.Lock()
	defer c.LockMutex.Unlock()
	if _, ok := c.Lock[userID]; !ok {
		c.Lock[userID] = &sync.Mutex{}
	}
	c.Lock[userID].Lock()
	defer c.Lock[userID].Unlock()
	if _, ok := c.UserConnections[userID]; !ok {
		c.UserConnections[userID] = conn
		return nil
	}
	return fmt.Errorf("connection with user id %s exists", userID)
}

/*
GetConn locks the mutex associated with a specific connection, before returning the connection itself.
After use, call the Release function to unlock the mutex.
*/
func (c *WsStore) GetConn(userID string) wsconnadapter.Adapter {
	c.LockMutex.Lock()
	defer c.LockMutex.Unlock()
	if _, ok := c.Lock[userID]; !ok {
		return nil
	}
	c.Lock[userID].Lock()
	if val, ok := c.UserConnections[userID]; ok {
		return val
	}
	return nil
}

// Release releases the lock for the specific connection
func (c *WsStore) Release(userID string) {
	c.LockMutex.Lock()
	defer c.LockMutex.Unlock()
	c.Lock[userID].Unlock() // There will always be a lock with userID for unlocking since deletion only occurs if we obtain the lock.
}

// DeleteConn deletes a user entry in the online users map
func (c *WsStore) DeleteConn(userID string) {
	c.LockMutex.Lock()
	defer c.LockMutex.Unlock()
	//We lock this for the final time before removing it, hence the absence of the unlock.
	if _, ok := c.Lock[userID]; ok {
		c.Lock[userID].Lock()
	}
	// Deletes are safe if key value pairs with supplied key do not exist
	go delete(c.UserConnections, userID)
	go delete(c.Lock, userID)
}
