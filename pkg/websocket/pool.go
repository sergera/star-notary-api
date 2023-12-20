package websocket

type PoolInterface interface {
	Start()
	RegisterConnection(ConnectionInterface)
	UnregisterConnection(ConnectionInterface)
	BroadcastJSONMessage(interface{})
}

type Pool struct {
	Register      chan ConnectionInterface
	Unregister    chan ConnectionInterface
	Connections   map[ConnectionInterface]bool
	BroadcastJSON chan interface{}
}

func NewPool() *Pool {
	p := &Pool{
		Register:      make(chan ConnectionInterface),
		Unregister:    make(chan ConnectionInterface),
		Connections:   make(map[ConnectionInterface]bool),
		BroadcastJSON: make(chan interface{}),
	}
	go p.Start()
	return p
}

func (pool *Pool) Start() {
	for {
		select {
		case conn := <-pool.Register:
			pool.Connections[conn] = true
		case conn := <-pool.Unregister:
			delete(pool.Connections, conn)
		case msg := <-pool.BroadcastJSON:
			for conn := range pool.Connections {
				go conn.WriteJSON(msg)
			}
		}
	}
}

func (pool *Pool) RegisterConnection(conn ConnectionInterface) {
	pool.Register <- conn
}

func (pool *Pool) UnregisterConnection(conn ConnectionInterface) {
	pool.Unregister <- conn
}

func (pool *Pool) BroadcastJSONMessage(msg interface{}) {
	pool.BroadcastJSON <- msg
}
