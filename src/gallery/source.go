package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/graphics-go/graphics"
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

func thumbnail(w http.ResponseWriter, fn string) {
	td := path.Join(path.Dir(fn), ".thumbnails")
	os.MkdirAll(td, 0700)
	tn := path.Join(td, path.Base(fn))
	f, err := os.Open(tn)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
		go func() {
			defer func() {
				if e := recover(); e != nil {
					fmt.Fprintln(os.Stderr, e.(error).Error())
				}
			}()
			f, err := os.Open(fn)
			assert(err)
			defer f.Close()
			g, err := os.Create(tn)
			assert(err)
			defer g.Close()
			var src image.Image
			switch strings.ToLower(path.Ext(fn)) {
			case ".png":
				src, err = png.Decode(f)
				assert(err)
				dst := image.NewRGBA(image.Rect(0, 0, 80, 80))
				graphics.Thumbnail(dst, src)
				assert(png.Encode(g, dst))
			default: //jpeg
				src, err = jpeg.Decode(f)
				assert(err)
				dst := image.NewRGBA(image.Rect(0, 0, 80, 80))
				graphics.Thumbnail(dst, src)
				assert(jpeg.Encode(g, dst, &jpeg.Options{jpeg.DefaultQuality}))
			}
		}()
		f, err = os.Open(fn)
		assert(err)
	}
	defer f.Close()
	switch strings.ToLower(path.Ext(fn)) {
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	default:
		w.Header().Set("Content-Type", "image/jpeg")
	}
	_, err = io.Copy(w, f)
	assert(err)
}

func source(w http.ResponseWriter, r *http.Request) {
	fn := path.Join(imgRoot, r.URL.Path[4:])
	st, err := os.Stat(fn)
	if err != nil || st.IsDir() {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
	if strings.Contains(r.URL.String(), "?thumbnail") {
		thumbnail(w, fn)
	} else {
		http.ServeFile(w, r, fn)
	}
}
