package main

import (
	"archive/zip"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/graphics-go/graphics"
)

func genThumbnail(fn string) bool {
	defer func() {
		if e := recover(); e != nil {
			fmt.Fprintln(os.Stderr, e.(error).Error())
		}
	}()
	tn := path.Join(path.Dir(fn), ".thumbnails", path.Base(fn))
	st, err := os.Stat(tn)
	if err == nil && !st.IsDir() && st.Size() > 0 {
		return false
	}
	g, err := os.Create(tn)
	assert(err)
	defer g.Close()
	f, err := os.Open(fn)
	assert(err)
	defer f.Close()
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
	return true
}

func GenerateResources() {
	zfn := path.Join(imgRoot, path.Base(imgRoot)+".zip")
	zf, err := os.Create(zfn)
	assert(err)
	defer zf.Close()
	zw := zip.NewWriter(zf)
	defer zw.Close()
	assert(os.MkdirAll(path.Join(imgRoot, ".thumbnails"), 0700))
	fns, _ := filepath.Glob(path.Join(imgRoot, "*"))
	fmt.Printf("Generating thumbnails and zipping %d resources...\n", len(fns))
	start := time.Now()
	for _, fn := range fns {
		switch strings.ToLower(path.Ext(fn)) {
		case ".jpg", ".jpeg", ".png":
			t := time.Now()
			if genThumbnail(fn) {
				fmt.Printf("  %s thumbnailed in %0.2f seconds\n", path.Base(fn), time.Since(t).Seconds())
			} else {
				fmt.Printf("  %s thumbnailing skipped (already exists)\n", path.Base(fn))
			}
			func() {
				t := time.Now()
				f, _ := os.Open(fn)
				defer func() {
					f.Close()
					fmt.Printf("  %s zipped in %0.2f seconds\n", path.Base(fn), time.Since(t).Seconds())
				}()
				info, _ := f.Stat()
				h, _ := zip.FileInfoHeader(info)
				h.Name = path.Base(fn)
				h.Method = zip.Store
				w, _ := zw.CreateHeader(h)
				io.Copy(w, f)
			}()
		}
	}
	fmt.Printf("resource preparation finished in %0.2f seconds.\n", time.Since(start).Seconds())
}
