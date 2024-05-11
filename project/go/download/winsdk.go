package download

import (
	"bytes"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	aflag "github.com/ricochhet/sdkstandalone/flag"
	"github.com/ricochhet/sdkstandalone/process"
	"github.com/tidwall/gjson"
)

func Getmsvcversion(aoutput string) (string, error) {
	dirPath := filepath.Join(aoutput, "VC", "Tools", "MSVC")
	fileList, err := filepath.Glob(filepath.Join(dirPath, "*"))
	if err != nil {
		return "", err
	}
	if len(fileList) == 0 {
		return "", fmt.Errorf("No files found in the directory:", dirPath)
	}

	return filepath.Base(fileList[0]), nil
}

func Getwinsdkversion(f *aflag.Flags) (string, error) {
	dirPath := filepath.Join(f.OUTPUT, "Windows Kits", "10", "bin")
	fileList, err := filepath.Glob(filepath.Join(dirPath, "*"))
	if err != nil {
		return "", err
	}
	if len(fileList) == 0 {
		return "", fmt.Errorf("No files found in the directory:", dirPath)
	}

	return filepath.Base(fileList[0]), nil
}

func Getwinsdk(f *aflag.Flags, json []gjson.Result, sdkPackages []string) error {
	var msi []string
	var cabs []string
	for _, item := range json {
		name := item.Get("fileName").String()
		url := item.Get("url").String()
		sha256 := item.Get("sha256").String()
		if slices.Contains(sdkPackages, name) {
			s := strings.TrimPrefix(name, "Installers\\")
			b, err := Downloadprogress(url, sha256, s, f.DOWNLOADS, s)
			if err != nil {
				fmt.Println("Error downloading Windows SDK package:", err)
				continue
			}

			msi = append(msi, filepath.Join(f.DOWNLOADS, s))
			cabs = append(cabs, Getmsicabs(b)...)
		}
	}

	for _, item := range json {
		name := item.Get("fileName").String()
		url := item.Get("url").String()
		sha256 := item.Get("sha256").String()
		if slices.Contains(cabs, strings.TrimPrefix(name, "Installers\\")) {
			s := strings.TrimPrefix(name, "Installers\\")
			_, err := Downloadprogress(url, sha256, s, f.DOWNLOADS, s)
			if err != nil {
				fmt.Println("Error downloading cab:", err)
				continue
			}
		}
	}

	for _, item := range msi {
		err := process.Exec("./rust-msiexec.exe", item, f.OUTPUT)
		if err != nil {
			return err
		}
	}

	return nil
}

func Getmsicabs(msi []byte) []string {
	var cabs []string
	index := 0
	for {
		foundIndex := bytes.Index(msi[index:], []byte(".cab"))
		if foundIndex < 0 {
			break
		}
		index += foundIndex + 4
		start := index - 36
		if start < 0 {
			start = 0
		}
		cabs = append(cabs, string(msi[start:index]))
	}
	return cabs
}
