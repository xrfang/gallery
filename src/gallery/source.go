package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path"
	"strings"
)

func sources(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			err := e.(error)
			if os.IsNotExist(err) {
				http.Error(w, "404 page not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}()
	dir, err := os.Open(imgRoot)
	assert(err)
	defer dir.Close()
	fis, err := dir.Readdir(65536)
	assert(err)
	var ps []string
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		switch strings.ToLower(path.Ext(fi.Name())) {
		case ".jpg", ".jpeg", ".png":
			ps = append(ps, fi.Name())
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ps)
}

func source(w http.ResponseWriter, r *http.Request) {
	fn := path.Join(imgRoot, r.URL.Path[4:])
	if strings.Contains(r.URL.String(), "?thumbnail") {
		fn = path.Join(path.Dir(fn), ".thumbnails", path.Base(fn))
	}
	http.ServeFile(w, r, fn)
}
