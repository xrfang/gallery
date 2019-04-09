package main

import (
	"net/http"
	"path"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.ServeFile(w, r, path.Join(webRoot, r.URL.Path))
		return
	}
	renderTemplate(w, "home.html", nil)
}
