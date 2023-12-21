package notifier

import (
	"net/http"
)

type MockStarNotifier struct {
	SubscribeCalled bool
	PublishCalled   bool
}

func (m *MockStarNotifier) Subscribe(w http.ResponseWriter, r *http.Request) {
	m.SubscribeCalled = true
}

func (m *MockStarNotifier) Publish(interface{}) {
	m.PublishCalled = true
}
