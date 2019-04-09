package main

import (
	"html/template"
	"net/http"
	"path"
	"path/filepath"
)

func renderTemplate(w http.ResponseWriter, tpl string, args interface{}) {
	defer func() {
		if e := recover(); e != nil {
			http.Error(w, e.(error).Error(), http.StatusInternalServerError)
		}
	}()
	helper := template.FuncMap{
		"ver": func() string {
			return "V" + _G_REVS + "." + _G_HASH
		},
	}
	tDir := path.Join(webRoot, "templates")
	t, err := template.New("body").Funcs(helper).ParseFiles(path.Join(tDir, tpl))
	assert(err)
	sfs, err := filepath.Glob(path.Join(tDir, "shared/*"))
	if len(sfs) > 0 {
		t, err = t.ParseFiles(sfs...)
		assert(err)
	}
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	assert(t.Execute(w, args))
}
