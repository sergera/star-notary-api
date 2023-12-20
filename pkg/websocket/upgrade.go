package websocket

import (
	"log"
	"net/http"
	"strings"

	"github.com/sergera/star-notary-backend/internal/conf"
	"nhooyr.io/websocket"
)

var wsOptions *websocket.AcceptOptions
var wsOptionsInitialized bool

func getWSOptions() *websocket.AcceptOptions {
	if !wsOptionsInitialized {
		wsOptions = &websocket.AcceptOptions{
			InsecureSkipVerify: true,
			OriginPatterns:     strings.Split(conf.ConfSingleton().CORSAllowedURLs, ","),
		}
		wsOptionsInitialized = true
	}
	return wsOptions
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := websocket.Accept(w, r, getWSOptions())
	if err != nil {
		log.Println(err)
		return ws, err
	}
	return ws, nil
}
