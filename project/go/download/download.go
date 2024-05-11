package download

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	aflag "github.com/ricochhet/sdkstandalone/flag"
)

var vstipExe = "vctip.exe"

func Createdirectories(f *aflag.Flags) (string, error) {
	err := os.MkdirAll(f.DOWNLOADS, 0700)
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(f.DOWNLOADS_CRTD, 0700)
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(f.DOWNLOADS_DIA, 0700)
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(f.OUTPUT, 0700)
	if err != nil {
		return "", err
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return wd, nil
}

func Removetelemetry(f *aflag.Flags) error {
	vctipX64 := filepath.Join(f.OUTPUT, "VC", "Tools", "MSVC", f.MSVC_VERSION_LOCAL, "bin", "Host"+f.HOST, f.TARGETX64, vstipExe)
	vctipX86 := filepath.Join(f.OUTPUT, "VC", "Tools", "MSVC", f.MSVC_VERSION_LOCAL, "bin", "Host"+f.HOST, f.TARGETX86, vstipExe)
	vctipARM := filepath.Join(f.OUTPUT, "VC", "Tools", "MSVC", f.MSVC_VERSION_LOCAL, "bin", "Host"+f.HOST, f.TARGETARM, vstipExe)
	vctipARM64 := filepath.Join(f.OUTPUT, "VC", "Tools", "MSVC", f.MSVC_VERSION_LOCAL, "bin", "Host"+f.HOST, f.TARGETARM64, vstipExe)

	err := os.Remove(vctipX64)
	if err != nil {
		return err
	}

	err = os.Remove(vctipX86)
	if err != nil {
		return err
	}

	if f.DOWNLOAD_ARM_TARGETS {
		err := os.Remove(vctipARM)
		if err != nil {
			return err
		}

		err = os.Remove(vctipARM64)
		if err != nil {
			return err
		}
	}

	return err
}

func Cleanhost(f *aflag.Flags) error {
	targets := []string{f.TARGETX64, f.TARGETX86, f.TARGETARM, f.TARGETARM64}
	for _, arch := range targets {
		if arch != f.HOST {
			err := os.RemoveAll(filepath.Join(f.OUTPUT, "VC", "Tools", "MSVC", f.MSVC_VERSION_LOCAL, "bin", "Host"+arch))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func Download(url string) ([]byte, error) {
	if url == "" {
		return nil, errors.New("url is empty")
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func Downloadprogress(url, check, name, apath, aname string) ([]byte, error) {
	if url == "" {
		return nil, errors.New("url is empty")
	}
	if apath == "" {
		return nil, errors.New("path is empty")
	}
	if name == "" || aname == "" {
		return nil, errors.New("name is empty")
	}

	fmt.Printf("%s ... DOWNLOADING\n", name)
	fpath := filepath.Join(apath, aname)
	err := os.MkdirAll(filepath.Dir(fpath), 0700)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(fpath); err == nil {
		data, err := os.ReadFile(fpath)
		if err == nil {
			hash := sha256.New()
			hash.Write(data)
			hashSum := hex.EncodeToString(hash.Sum(nil))
			if strings.ToLower(check) == hashSum {
				fmt.Printf("%s ... OK\n", name)
				return data, nil
			}
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	f, err := os.Create(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	hash := sha256.New()
	buf := make([]byte, 1<<20)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}
		if _, err := f.Write(buf[:n]); err != nil {
			return nil, err
		}
		if _, err := hash.Write(buf[:n]); err != nil {
			return nil, err
		}
	}
	data, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	hashsum := hex.EncodeToString(hash.Sum(nil))
	if strings.ToLower(check) != hashsum {
		return nil, fmt.Errorf("hash mismatch for %s", name)
	}
	return data, nil
}
