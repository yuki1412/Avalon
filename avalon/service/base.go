package service

import (
	"net/http"
	"text/template"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/index.html"))
	res := map[string]interface{}{}
	tmpl.Execute(w, res)
	// fmt.Fprintf(w, "Welcome to the home page!")
}
