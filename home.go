package main

import "net/http"

func showHome(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "home.gohtml", nil)
}

func showDefault(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "default.gohtml", nil)
}
