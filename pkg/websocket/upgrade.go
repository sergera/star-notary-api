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
		conf, err := conf.ConfSingleton()
		var originPatterns []string
		if err == nil {
			originPatterns = strings.Split(conf.CORSAllowedURLs, ",")
		}
		wsOptions = &websocket.AcceptOptions{
			InsecureSkipVerify: true,
			OriginPatterns:     originPatterns,
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
