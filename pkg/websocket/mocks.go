package websocket

import (
	"context"
	"sync"

	"nhooyr.io/websocket"
)

type MockPool struct {
	RegisterCalled      bool
	UnregisterCalled    bool
	BroadcastJSONCalled bool
	BroadcastMsg        interface{}
}

func (m *MockPool) Start() {}

func (m *MockPool) RegisterConnection(ConnectionInterface) {
	m.RegisterCalled = true
}

func (m *MockPool) UnregisterConnection(ConnectionInterface) {
	m.UnregisterCalled = true
}

func (m *MockPool) BroadcastJSONMessage(msg interface{}) {
	m.BroadcastJSONCalled = true
	m.BroadcastMsg = msg
}

type MockConnection struct {
	WriteJSONCalled bool
	ReceivedMessage interface{}
	wg              *sync.WaitGroup
}

func (m *MockConnection) WriteText(body []byte) error {
	m.ReceivedMessage = body
	return nil
}

func (m *MockConnection) WriteBinary(body []byte) error {
	m.ReceivedMessage = body
	return nil
}

func (m *MockConnection) WriteJSON(body interface{}) error {
	m.WriteJSONCalled = true
	m.ReceivedMessage = body
	if m.wg != nil {
		m.wg.Done()
	}
	return nil
}

func (m *MockConnection) Close(status websocket.StatusCode, reason string) {}

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
