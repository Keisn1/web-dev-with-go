package views

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func Parse(filepath string) (Template, error) {
	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}
	return Template{hmtlTpl: tpl}, nil
}

type Template struct {
	hmtlTpl *template.Template
}

func (t *Template) Execute(w http.ResponseWriter, data any) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	err := t.hmtlTpl.Execute(w, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
		return
	}
}
