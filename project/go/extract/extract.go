package extract

import (
	"archive/zip"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func Extractpackage(fpath, destpath string) {
	if !strings.HasSuffix(fpath, ".vsix") {
		fmt.Println("File is not a vsix file.")
		return
	}

	zr, err := zip.OpenReader(fpath)
	if err != nil {
		fmt.Println("Error opening vsix file:", err)
		return
	}

	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			fmt.Println("Error opening vsix file for extraction:", err)
			continue
		}
		defer rc.Close()
		if !strings.HasPrefix(f.Name, "Contents") {
			continue
		}

		name, err := url.QueryUnescape(strings.TrimPrefix(f.Name, "Contents/"))
		if err != nil {
			fmt.Println("Error decoding name:", err)
			return
		}
		dest := filepath.Join(destpath, name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(dest, os.ModePerm)
		} else {
			os.MkdirAll(filepath.Dir(dest), os.ModePerm)
			w, err := os.Create(dest)
			if err != nil {
				fmt.Println("Error creating destination file:", err)
				continue
			}
			if _, err := io.Copy(w, rc); err != nil {
				fmt.Println("Error copying file:", err)
				w.Close()
				continue
			}
			w.Close()
		}
	}
	zr.Close()
}
