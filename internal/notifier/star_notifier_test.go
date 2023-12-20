package notifier

import (
	"net/http/httptest"
	"testing"

	"github.com/sergera/star-notary-backend/pkg/websocket"
)

func TestStarNotifier_Subscribe(t *testing.T) {
	mockPool := &websocket.MockPool{}
	notifier := &StarNotifier{pool: mockPool}

	// mock http.ResponseWriter and http.Request
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/subscribe", nil)

	// call Subscribe
	notifier.Subscribe(w, r)

	// verify that RegisterConnection was called
	if !mockPool.RegisterCalled {
		t.Errorf("Subscribe did not register the connection")
	}
}

func TestStarNotifier_Publish(t *testing.T) {
	mockPool := &websocket.MockPool{}
	notifier := &StarNotifier{pool: mockPool}

	testMsg := "test message"
	notifier.Publish(testMsg)

	// verify that BroadcastJSONMessage was called with the correct message
	if !mockPool.BroadcastJSONCalled || mockPool.BroadcastMsg != testMsg {
		t.Errorf("Publish did not broadcast the correct message")
	}
}
