package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	res "github.com/xrfang/go-res"
)

var imgRoot, webRoot string

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var uri string
	ver := flag.Bool("version", false, "show version info")
	port := flag.String("port", "8080", "service port")
	flag.StringVar(&uri, "uri", "/", "relative path to the gallery")
	flag.StringVar(&webRoot, "webroot", "/tmp/gallery/webroot", "")
	flag.StringVar(&imgRoot, "imgs", "", "directory for images")
	pkg := flag.String("pack", "", "pack resources under directory")
	flag.Parse()
	if *ver {
		fmt.Println(verinfo())
		return
	}
	if *pkg != "" {
		assert(res.Pack(*pkg))
		fmt.Printf("resources under '%s' packed.\n", *pkg)
		return
	}
	st, err := os.Stat(imgRoot)
	if err != nil || !st.IsDir() {
		fmt.Println("invalid or missing image directory (-imgs)")
		os.Exit(1)
	}
	assert(res.Extract(webRoot, res.OverwriteIfNewer))
	if !strings.HasPrefix(uri, "/") {
		uri = "/" + uri
	}
	http.HandleFunc(uri, home)
	http.HandleFunc(uri+"imgs", sources)
	http.HandleFunc(uri+"img/", source)
	svr := http.Server{
		Addr:         ":" + *port,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
	assert(svr.ListenAndServe())
}
