package extract

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var errNotVsixFile = errors.New("file is not a vsix file")

func Vsix(fpath, destpath string) error {
	if !strings.HasSuffix(fpath, ".vsix") {
		return errNotVsixFile
	}

	return extractZip(fpath, destpath)
}

func extractZip(fpath, destpath string) error {
	copybytes := 1024
	zipread, err := zip.OpenReader(fpath)
	if err != nil { //nolint:wsl // gofumpt conflict
		return err
	}

	for _, file := range zipread.File {
		readclose, err := file.Open()
		if err != nil {
			fmt.Println("Error opening zip file for extraction:", err)
			continue
		}
		defer readclose.Close()

		if !strings.HasPrefix(file.Name, "Contents") {
			continue
		}

		name, err := url.QueryUnescape(strings.TrimPrefix(file.Name, "Contents/"))
		if err != nil {
			fmt.Println("Error decoding name:", err)
			return err
		}

		dest := filepath.Join(destpath, name)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(dest, os.ModePerm); err != nil {
				continue
			}
		} else if err := destcopy(dest, readclose, int64(copybytes)); err != nil {
			return err
		}
	}

	zipread.Close()

	return nil
}

func destcopy(apath string, src io.Reader, copybytes int64) error {
	if err := os.MkdirAll(filepath.Dir(apath), os.ModePerm); err != nil {
		return err
	}

	dst, err := os.Create(apath)
	if err != nil {
		return err
	}

	for {
		_, err := io.CopyN(dst, src, copybytes)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return err
		}
	}

	dst.Close()

	return nil
}
