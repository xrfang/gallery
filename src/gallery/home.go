package main

import (
	"net/http"
	"os"
	"path"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.ServeFile(w, r, path.Join(webRoot, r.URL.Path))
		return
	}
	zfn := path.Join(imgRoot, path.Base(imgRoot)+".zip")
	st, _ := os.Stat(zfn)
	renderTemplate(w, "home.html", struct {
		Title string
		Zip   string
		Total int
	}{
		Title: galleryTitle,
		Zip:   uri + "img/" + path.Base(zfn),
		Total: int(st.Size() / 1024 / 1024),
	})
}
