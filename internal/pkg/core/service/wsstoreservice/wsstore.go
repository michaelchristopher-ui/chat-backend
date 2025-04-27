package wsstoreservice

import (
	"fmt"
	"sync"
	"websocket_client/internal/pkg/core/adapter/wsconnadapter"
	"websocket_client/internal/pkg/core/adapter/wsstoreadapter"
)

type NewWsStoreServiceReq struct {
	UserConnections map[string]wsconnadapter.Adapter
	ConnLock        map[string]sync.Locker
	LockMutex       sync.Locker
}

type WsStore struct {
	userConnections map[string]wsconnadapter.Adapter //In memory store for the connections
	connLock        map[string]sync.Locker           //Mutex for the connections
	lockMutex       sync.Locker                      //Mutex for the maps
}

func NewWSStoreService(req NewWsStoreServiceReq) wsstoreadapter.Adapter {
	return &WsStore{
		userConnections: req.UserConnections,
		connLock:        req.ConnLock,
		lockMutex:       req.LockMutex,
	}
}

// AddConn adds a user connection key value entry in the online users map provided that it is able to obtain all required locks
func (c *WsStore) AddConn(conn wsconnadapter.Adapter, userID string) error {
	c.lockMutex.Lock()
	defer c.lockMutex.Unlock()
	if _, ok := c.connLock[userID]; !ok {
		c.connLock[userID] = &sync.Mutex{}
	}
	c.connLock[userID].Lock()
	defer c.connLock[userID].Unlock()
	if _, ok := c.userConnections[userID]; !ok {
		c.userConnections[userID] = conn
		return nil
	}
	return fmt.Errorf(ErrConnExistString, userID)
}

/*
GetConn locks the mutex associated with a specific connection, before returning the connection itself.
After use, call the Release function to unlock the mutex.
*/
func (c *WsStore) GetConn(userID string) wsconnadapter.Adapter {
	c.lockMutex.Lock()
	defer c.lockMutex.Unlock()
	if _, ok := c.connLock[userID]; !ok {
		return nil
	}
	c.connLock[userID].Lock()
	if val, ok := c.userConnections[userID]; ok {
		return val
	}
	return nil
}

// Release releases the lock for the specific connection
func (c *WsStore) Release(userID string) {
	c.lockMutex.Lock()
	defer c.lockMutex.Unlock()
	c.connLock[userID].Unlock() // There will always be a lock with userID for unlocking since deletion only occurs if we obtain the lock.
}

// DeleteConn deletes a user entry in the online users map
func (c *WsStore) DeleteConn(userID string) {
	c.lockMutex.Lock()
	defer c.lockMutex.Unlock()
	//We lock this for the final time before removing it, hence the absence of the unlock.
	if _, ok := c.connLock[userID]; ok {
		c.connLock[userID].Lock()
	}
	// Deletes are safe if key value pairs with supplied key do not exist
	go delete(c.userConnections, userID)
	go delete(c.connLock, userID)
}
