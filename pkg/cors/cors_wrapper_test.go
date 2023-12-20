package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestNewCors checks if CORS struct is initialized correctly
func TestNewCors(t *testing.T) {
	urls := []string{"http://example.com", "http://test.com"}
	verbs := []HTTPVerb{Get, Post}

	cors := NewCors(urls, verbs)
	expectedURLs := "http://example.com, http://test.com"
	expectedVerbs := "Get, Post"

	if cors.allowedURLPatterns != expectedURLs {
		t.Errorf("Expected URL patterns %s, got %s", expectedURLs, cors.allowedURLPatterns)
	}
	if cors.allowedVerbs != expectedVerbs {
		t.Errorf("Expected verbs %s, got %s", expectedVerbs, cors.allowedVerbs)
	}
}

// TestWrapHandlerFunc checks if the middleware modifies the response correctly
func TestWrapHandlerFunc(t *testing.T) {
	urls := []string{"*"}
	verbs := []HTTPVerb{Get, Post}
	cors := NewCors(urls, verbs)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrappedHandler := cors.WrapHandlerFunc(testHandler)
	testServer := httptest.NewServer(wrappedHandler)
	defer testServer.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", testServer.URL, nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("Failed to execute request")
	}

	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Expected Access-Control-Allow-Origin header to be set")
	}
}

// TestWrapHandler checks if the middleware modifies the response correctly for http.Handler
func TestWrapHandler(t *testing.T) {
	urls := []string{"*"}
	verbs := []HTTPVerb{Get, Post}
	cors := NewCors(urls, verbs)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrappedHandler := cors.WrapHandler(testHandler)
	testServer := httptest.NewServer(wrappedHandler)
	defer testServer.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", testServer.URL, nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("Failed to execute request")
	}

	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Expected Access-Control-Allow-Origin header to be set")
	}

	// Testing preflight request
	preflightReq, _ := http.NewRequest("OPTIONS", testServer.URL, nil)
	preflightResp, err := client.Do(preflightReq)
	if err != nil {
		t.Fatal("Failed to execute preflight request")
	}

	if preflightResp.Header.Get("Access-Control-Allow-Methods") != "Get, Post" {
		t.Error("Expected Access-Control-Allow-Methods header to be set for preflight request")
	}
}
