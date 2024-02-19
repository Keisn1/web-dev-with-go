package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/keisn1/lenslocked/controllers"
	"github.com/keisn1/lenslocked/views"
	// "github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	// https://devdocs.io/go/net/http/index#HandleFunc
	// HandleFunc registers the handler function for the given pattern in the
	// DefaultServeMux. The documentation for ServeMux explains how patterns are
	// matched.
	// http.HandleFunc("/", pathHandler)

	// ListenAndServe listens on the TCP network address addr and then calls Serve
	// with handler to handle requests on incoming connections. Accepted connections
	// are configured to enable TCP keep-alives.

	homeTpl := views.Must(views.Parse(filepath.Join("templates", "home.html")))
	r.Get("/", controllers.StaticHandler(homeTpl))
	contactTpl := views.Must(views.Parse(filepath.Join("templates", "contact.html")))
	r.Route("/contact", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/", controllers.StaticHandler(contactTpl))
	})

	faqTpl := views.Must(views.Parse(filepath.Join("templates", "faq.html")))
	r.Get("/faq", controllers.StaticHandler(faqTpl))

	r.Get("/gallery/{galleryID}", galleryHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not foundicilious", http.StatusNotFound)
	})
	fmt.Println("Starting the server on: 3000...")
	http.ListenAndServe(":3000", r)
}

func executeTemplate(w http.ResponseWriter, filepath string) {
	tpl, err := views.Parse(filepath)
	if err != nil {
		log.Printf("parsing template: %v", err)
		return
	}
	tpl.Execute(w, nil)
}

func galleryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<h1>You requested gallery number %s</h1>", chi.URLParam(r, "galleryID"))
	return
}
