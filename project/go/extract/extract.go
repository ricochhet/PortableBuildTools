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

func Extractpackage(apath, aoutput string) {
	if !strings.HasSuffix(apath, ".vsix") {
		fmt.Println("File is not a zip file.")
		return
	}

	zr, err := zip.OpenReader(apath)
	if err != nil {
		fmt.Println("Error opening zip file:", err)
		return
	}

	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			fmt.Println("Error opening zip file for extraction:", err)
			continue
		}
		defer rc.Close()
		if !strings.HasPrefix(f.Name, "Contents") {
			continue
		}

		name, found := strings.CutPrefix(f.Name, "Contents/")
		if !found {
			break
		}

		decodedName, err := url.QueryUnescape(name)
		if err != nil {
			fmt.Println("Error decoding name:", err)
			return
		}
		destPath := filepath.Join(aoutput, decodedName)
		if f.FileInfo().IsDir() {
			os.MkdirAll(destPath, os.ModePerm)
		} else {
			os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
			w, err := os.Create(destPath)
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
