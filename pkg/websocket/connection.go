package websocket

import (
	"context"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const (
	// time allowed to write a message to the peer
	writeWait = 10 * time.Second
)

type ConnectionInterface interface {
	WriteText(body []byte) error
	WriteBinary(body []byte) error
	WriteJSON(body interface{}) error
	Close(status websocket.StatusCode, reason string)
}

type WebsocketConnWrapperInterface interface {
	Subprotocol() string
	Ping(ctx context.Context) error
	Read(ctx context.Context) (websocket.MessageType, []byte, error)
	Write(ctx context.Context, messageType websocket.MessageType, payload []byte) error
	Close(status websocket.StatusCode, reason string) error
	WriteJSON(ctx context.Context, body interface{}) error
}

type websocketConnWrapper struct {
	inner *websocket.Conn
}

func NewWebsocketConnWrapper(inner *websocket.Conn) *websocketConnWrapper {
	return &websocketConnWrapper{inner}
}

func (w *websocketConnWrapper) Subprotocol() string {
	return w.inner.Subprotocol()
}

func (w *websocketConnWrapper) Ping(ctx context.Context) error {
	return w.inner.Ping(ctx)
}

func (w *websocketConnWrapper) Read(ctx context.Context) (websocket.MessageType, []byte, error) {
	return w.inner.Read(ctx)
}

func (w *websocketConnWrapper) Write(ctx context.Context, messageType websocket.MessageType, payload []byte) error {
	return w.inner.Write(ctx, messageType, payload)
}

func (w *websocketConnWrapper) Close(status websocket.StatusCode, reason string) error {
	return w.inner.Close(status, reason)
}

func (w *websocketConnWrapper) WriteJSON(ctx context.Context, body interface{}) error {
	return wsjson.Write(ctx, w.inner, body)
}

type Connection struct {
	inner WebsocketConnWrapperInterface
	pool  PoolInterface
}

type Message struct {
	Type websocket.MessageType `json:"type"`
	Body string                `json:"body"`
}

func NewConnection(conn WebsocketConnWrapperInterface, pool PoolInterface) *Connection {
	return &Connection{conn, pool}
}

func (c *Connection) WriteText(body []byte) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(writeWait))
	defer cancel()
	err := c.inner.Write(ctx, websocket.MessageText, body)
	if err != nil {
		c.Close(websocket.StatusAbnormalClosure, "error writing text message")
		return err
	}
	return nil
}

func (c *Connection) WriteBinary(body []byte) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(writeWait))
	defer cancel()
	err := c.inner.Write(ctx, websocket.MessageBinary, body)
	if err != nil {
		c.Close(websocket.StatusAbnormalClosure, "error writing binary message")
		return err
	}
	return nil
}

func (c *Connection) WriteJSON(body interface{}) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(writeWait))
	defer cancel()
	err := c.inner.WriteJSON(ctx, body)
	if err != nil {
		c.Close(websocket.StatusAbnormalClosure, "error writing JSON message")
		return err
	}
	return nil
}

func (c *Connection) Close(status websocket.StatusCode, reason string) {
	c.pool.UnregisterConnection(c)
	c.inner.Close(status, reason)
}
