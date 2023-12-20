package websocket

import (
	"context"
	"testing"

	"nhooyr.io/websocket"
)

type MockPool struct {
	UnregisterCalled bool
}

func (m *MockPool) Start() {}

func (m *MockPool) UnregisterConnection(c *Connection) {
	m.UnregisterCalled = true
}

type MockWebsocketConn struct {
	WriteJSONCalled bool
	WriteJSONErr    error
	CloseCalled     bool
}

func (m *MockWebsocketConn) Subprotocol() string {
	return ""
}

func (m *MockWebsocketConn) Ping(ctx context.Context) error {
	return nil
}

func (m *MockWebsocketConn) Read(ctx context.Context) (websocket.MessageType, []byte, error) {
	return websocket.MessageText, []byte{}, nil
}

func (m *MockWebsocketConn) Write(ctx context.Context, messageType websocket.MessageType, p []byte) error {
	return nil
}

func (m *MockWebsocketConn) Close(statusCode websocket.StatusCode, reason string) error {
	m.CloseCalled = true
	return nil
}

func (m *MockWebsocketConn) WriteJSON(ctx context.Context, v interface{}) error {
	m.WriteJSONCalled = true
	return m.WriteJSONErr
}

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
