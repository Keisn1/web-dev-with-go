package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	// "github.com/go-chi/chi/v5/middleware"
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

	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	r.Get("/", homeHandler)
	r.Route("/contact", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/", contactHandler)
	})
	r.Get("/faq", faqHandler)
	r.Get("/gallery/{galleryID}", galleryHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not foundicilious", http.StatusNotFound)
	})
	fmt.Println("Starting the server on: 3000...")
	http.ListenAndServe(":3000", r)
}

func executeTemplate(w http.ResponseWriter, filepath string) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
		return
	}

	err = tpl.Execute(w, nil)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "contact.html")
	executeTemplate(w, tplPath)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "faq.html")
	executeTemplate(w, tplPath)

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "home.html")
	executeTemplate(w, tplPath)
}

func galleryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<h1>You requested gallery number %s</h1>", chi.URLParam(r, "galleryID"))
	return
}
