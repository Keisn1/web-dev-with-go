package main

import (
	"fmt"
	"net/http"
)

type Router struct{}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	case "/faq":
		faqHandler(w, r)
	default:
		http.Error(w, "Page not found", http.StatusNotFound)
	}
}

func main() {
	// https://devdocs.io/go/net/http/index#HandleFunc
	// HandleFunc registers the handler function for the given pattern in the
	// DefaultServeMux. The documentation for ServeMux explains how patterns are
	// matched.
	// http.HandleFunc("/", pathHandler)

	// ListenAndServe listens on the TCP network address addr and then calls Serve
	// with handler to handle requests on incoming connections. Accepted connections
	// are configured to enable TCP keep-alives.

	var router Router
	fmt.Println("Starting the server on: 3000...")
	http.ListenAndServe(":3000", router)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Contact Page</h1><p>To get in touch, email me at <a href=\"mailto:jon@calho")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>FAQ Page</h1><p>Do you have questions?</p>")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<h1>Welcome to my awesome site</h1>")
}
