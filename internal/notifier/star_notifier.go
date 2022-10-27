package notifier

import (
	"context"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sergera/star-notary-backend/internal/conf"
	"github.com/sergera/star-notary-backend/internal/domain"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var once sync.Once
var instance *StarNotifier

type StarNotifier struct {
	wsOptions websocket.AcceptOptions
	mu        sync.Mutex
	wg        *sync.WaitGroup
	queue     []domain.StarModel
	ticker    *time.Ticker
}

func StarNotifierSingleton() *StarNotifier {
	once.Do(func() {
		conf := conf.ConfSingleton()
		var n *StarNotifier = &StarNotifier{
			websocket.AcceptOptions{
				InsecureSkipVerify: true,
				OriginPatterns:     strings.Split(conf.CORSAllowedURLs, ","),
			},
			sync.Mutex{},
			&sync.WaitGroup{},
			[]domain.StarModel{},
			time.NewTicker(time.Second),
		}
		instance = n
	})
	go instance.cleanQueue()
	return instance
}

func (n *StarNotifier) PushStars(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &n.wsOptions)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "oops, something went wrong")

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	ctx = c.CloseRead(ctx)

	for {
		select {
		case <-ctx.Done():
			c.Close(websocket.StatusNormalClosure, "")
			return
		case <-n.ticker.C:
			if len(n.queue) > 0 {
				n.wg.Add(1)
				if err := wsjson.Write(ctx, c, n.queue); err != nil {
					log.Println("error sending websocket message: ", err.Error())
				}
				n.wg.Done()
			}
		}
	}
}

func (n *StarNotifier) cleanQueue() {
	for range n.ticker.C {
		time.Sleep(500 * time.Millisecond)
		n.wg.Wait()
		n.mu.Lock()
		n.queue = []domain.StarModel{}
		n.mu.Unlock()
	}
}

func (n *StarNotifier) AppendStar(m domain.StarModel) {
	n.mu.Lock()
	n.queue = append(n.queue, m)
	n.mu.Unlock()
}
