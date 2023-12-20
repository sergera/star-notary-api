package websocket

import (
	"sync"
	"testing"
	"time"
)

func TestPool_RegisterUnregister(t *testing.T) {
	pool := NewPool()

	mockConn1 := &MockConnection{wg: new(sync.WaitGroup)}
	mockConn2 := &MockConnection{wg: new(sync.WaitGroup)}

	// setup WaitGroup to wait for 2 WriteJSON calls
	wg := new(sync.WaitGroup)
	wg.Add(2)
	mockConn1.wg = wg
	mockConn2.wg = wg

	// register connections
	pool.Register <- mockConn1
	pool.Register <- mockConn2

	// test broadcasting
	testMessage := "test message"
	pool.BroadcastJSON <- testMessage

	// wait for both WriteJSON calls to complete
	wg.Wait()

	// check if WriteJSON was called on both connections
	if !mockConn1.WriteJSONCalled || !mockConn2.WriteJSONCalled {
		t.Errorf("Broadcast message was not received by all connections")
	}

	pool.Unregister <- mockConn1

	// polling mechanism to wait for unregistration
	// TODO: replace with a channel
	const maxRetries = 100
	var retries int
	for retries < maxRetries {
		if _, exists := pool.Connections[mockConn1]; !exists {
			break // connection is unregistered successfully
		}
		time.Sleep(1 * time.Millisecond) // brief sleep between retries
		retries++
	}

	// check if the connection is removed from the pool
	if _, exists := pool.Connections[mockConn1]; exists {
		t.Errorf("Connection was not unregistered correctly")
	}
}

func TestPool_BroadcastsToAllConnections(t *testing.T) {
	var wg sync.WaitGroup
	pool := NewPool()

	// create and register mock connections
	mockConn1 := &MockConnection{wg: &wg}
	mockConn2 := &MockConnection{wg: &wg}

	wg.Add(2) // expect 2 calls to WriteJSON
	pool.Register <- mockConn1
	pool.Register <- mockConn2

	// broadcast a test message
	testMessage := "test message"
	pool.BroadcastJSON <- testMessage

	wg.Wait() // wait for all calls to WriteJSON to finish

	// verify that each connection received the broadcast message
	if !mockConn1.WriteJSONCalled || mockConn1.ReceivedMessage != testMessage {
		t.Errorf("Connection 1 did not receive the correct broadcast message")
	}
	if !mockConn2.WriteJSONCalled || mockConn2.ReceivedMessage != testMessage {
		t.Errorf("Connection 2 did not receive the correct broadcast message")
	}
}
