package download

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	aflag "github.com/ricochhet/sdkstandalone/flag"
	"github.com/tidwall/gjson"
)

var installerPrefix = "Installers\\"

func Getmsvcversion(destpath string) (string, error) {
	msvcbin := filepath.Join(destpath, "VC", "Tools", "MSVC")
	binversions, err := filepath.Glob(filepath.Join(msvcbin, "*"))
	if err != nil {
		return "", err
	}
	if len(binversions) == 0 {
		return "", errors.New("No files found in the directory:" + msvcbin)
	}

	return filepath.Base(binversions[0]), nil
}

func Getwinsdkversion(f *aflag.Flags) (string, error) {
	winsdkbin := filepath.Join(f.OUTPUT, "Windows Kits", "10", "bin")
	binversions, err := filepath.Glob(filepath.Join(winsdkbin, "*"))
	if err != nil {
		return "", err
	}
	if len(binversions) == 0 {
		return "", errors.New("No files found in the directory:" + winsdkbin)
	}

	return filepath.Base(binversions[0]), nil
}

func Getwinsdk(f *aflag.Flags, packages []gjson.Result, sdkPackages []string) error {
	installers := []string{}
	cabinets := []string{}
	for _, pkg := range packages {
		name := pkg.Get("fileName").String()
		url := pkg.Get("url").String()
		sha256 := pkg.Get("sha256").String()
		if slices.Contains(sdkPackages, name) {
			filename := strings.TrimPrefix(name, installerPrefix)
			installer, err := Downloadprogress(url, sha256, filename, f.DOWNLOADS, filename)
			if err != nil {
				fmt.Println("Error downloading Windows SDK package:", err)
				continue
			}

			installers = append(installers, filepath.Join(f.DOWNLOADS, filename))
			cabinets = append(cabinets, Getmsicabinets(installer)...)
		}
	}

	for _, pkg := range packages {
		name := pkg.Get("fileName").String()
		url := pkg.Get("url").String()
		sha256 := pkg.Get("sha256").String()
		if slices.Contains(cabinets, strings.TrimPrefix(name, installerPrefix)) {
			filename := strings.TrimPrefix(name, installerPrefix)
			_, err := Downloadprogress(url, sha256, filename, f.DOWNLOADS, filename)
			if err != nil {
				fmt.Println("Error downloading cab:", err)
				continue
			}
		}
	}

	for _, installer := range installers {
		err := Rustmsiexec(f, "./rust-msiexec.exe", installer, f.OUTPUT)
		if err != nil {
			return err
		}
	}

	return nil
}

func Getmsicabinets(installerMsi []byte) []string {
	cabs := []string{}
	index := 0
	for {
		foundIndex := bytes.Index(installerMsi[index:], []byte(".cab"))
		if foundIndex < 0 {
			break
		}
		index += foundIndex + 4
		start := index - 36
		if start < 0 {
			start = 0
		}
		cabs = append(cabs, string(installerMsi[start:index]))
	}
	return cabs
}
