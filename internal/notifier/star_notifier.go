package notifier

import (
	"log"
	"net/http"
	"sync"

	"github.com/sergera/star-notary-backend/pkg/websocket"
)

var once sync.Once
var instance *StarNotifier

type StarNotifier struct {
	pool *websocket.Pool
}

func StarNotifierSingleton() *StarNotifier {
	once.Do(func() {
		pool := websocket.NewPool()

		instance = &StarNotifier{pool}
	})
	return instance
}

func (n *StarNotifier) Subscribe(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r)
	if err != nil {
		log.Println(err)
	}

	n.pool.Register <- websocket.NewConnection(websocket.NewWebsocketConnWrapper(ws), n.pool)
}

func (n *StarNotifier) Publish(msg interface{}) {
	n.pool.BroadcastJSON <- msg
}
