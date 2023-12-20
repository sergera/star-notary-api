package websocket

import (
	"testing"

	"nhooyr.io/websocket"
)

func TestConnection_WriteAndClose(t *testing.T) {
	mockPool := &MockPool{}

	mockWebSocketConn := &MockWebsocketConn{}

	conn := NewConnection(mockWebSocketConn, mockPool)

	// test writing JSON
	err := conn.WriteJSON("test message")
	if err != nil {
		t.Errorf("WriteJSON failed: %v", err)
	}

	// test closing the connection
	conn.Close(websocket.StatusNormalClosure, "test reason")

	// check if UnregisterConnection was called in the pool
	if !mockPool.UnregisterCalled {
		t.Errorf("Close did not unregister the connection from the pool")
	}
}
