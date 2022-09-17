package cors

import (
	"net/http"
	"strings"

	"github.com/sergera/star-notary-backend/pkg/slc"
)

type HTTPVerb int8

const (
	Get HTTPVerb = iota
	Head
	Post
	Put
	Delete
	Connect
	Options
	Trace
	Patch
)

func (v HTTPVerb) String() string {
	switch v {
	case Get:
		return "Get"
	case Head:
		return "Head"
	case Post:
		return "Post"
	case Put:
		return "Put"
	case Delete:
		return "Delete"
	case Options:
		return "Options"
	case Trace:
		return "Trace"
	case Patch:
		return "Patch"
	default:
		return "Unknown"
	}
}

type Cors struct {
	allowedURLPatterns string
	allowedVerbs       string
}

func NewCors(urls []string, verbs []HTTPVerb) *Cors {
	urlsString := strings.Join(urls, ", ")
	verbsString := strings.Join(slc.Map(verbs, func(verb HTTPVerb) string {
		return verb.String()
	}), ", ")

	return &Cors{
		allowedURLPatterns: urlsString,
		allowedVerbs:       verbsString,
	}
}

func (c *Cors) handlePreflight(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Accept")
	w.Header().Add("Access-Control-Allow-Methods", c.allowedVerbs)
}

func (c *Cors) WrapHandlerFunc(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", c.allowedURLPatterns)
		if r.Method == "OPTIONS" {
			c.handlePreflight(w)
		} else {
			h.ServeHTTP(w, r)
		}
	}
}

func (c *Cors) WrapHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", c.allowedURLPatterns)
		if r.Method == "OPTIONS" {
			c.handlePreflight(w)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}
