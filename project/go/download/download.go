package download

import (
	"context"
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
	"github.com/ricochhet/sdkstandalone/process"
)

var (
	errURLEmpty  = errors.New("url is empty")
	errPathEmpty = errors.New("path is empty")
	errNameEmpty = errors.New("name is empty")
)

func CreateDirectories(flags *aflag.Flags) (string, error) {
	directories := []string{flags.Downloads, flags.DownloadsCRTD, flags.DownloadsDIA, flags.Output}
	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return "", err
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return wd, nil
}

func RemoveVCTipsTelemetry(flags *aflag.Flags) error {
	vctipExe := "vctip.exe"
	paths := []string{
		filepath.Join(flags.Output, "VC", "Tools", "MSVC", flags.MsvcVerLocal, "bin", "Host"+flags.Host, flags.Targetx64, vctipExe),
		filepath.Join(flags.Output, "VC", "Tools", "MSVC", flags.MsvcVerLocal, "bin", "Host"+flags.Host, flags.Targetx86, vctipExe),
	}

	if flags.DownloadARMTargets {
		paths = append(paths,
			filepath.Join(flags.Output, "VC", "Tools", "MSVC", flags.MsvcVerLocal, "bin", "Host"+flags.Host, flags.Targetarm, vctipExe),
			filepath.Join(flags.Output, "VC", "Tools", "MSVC", flags.MsvcVerLocal, "bin", "Host"+flags.Host, flags.Targetarm64, vctipExe),
		)
	}

	for _, path := range paths {
		err := os.Remove(path)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	return nil
}

func CleanHostDirectory(flags *aflag.Flags) error {
	targets := []string{flags.Targetx64, flags.Targetx86, flags.Targetarm, flags.Targetarm64}
	for _, arch := range targets {
		if arch != flags.Host {
			err := os.RemoveAll(filepath.Join(flags.Output, "VC", "Tools", "MSVC", flags.MsvcVerLocal, "bin", "Host"+arch))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func Download(url string) ([]byte, error) {
	if url == "" {
		return nil, errURLEmpty
	}

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
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

func File(url, check, name, apath, aname string) ([]byte, error) {
	if err := validateDownloadParams(url, apath, name, aname); err != nil {
		return nil, err
	}

	fpath := filepath.Join(apath, aname)
	if err := os.MkdirAll(filepath.Dir(fpath), 0o700); err != nil {
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

	fmt.Printf("%s ... DOWNLOADING\n", name)

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	file, err := os.Create(fpath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return hashMatch(resp, file, fpath, check, name)
}

func validateDownloadParams(url, apath, name, aname string) error {
	if url == "" {
		return errURLEmpty
	}

	if apath == "" {
		return errPathEmpty
	}

	if name == "" || aname == "" {
		return errNameEmpty
	}

	return nil
}

func hashMatch(resp *http.Response, flags *os.File, fpath, check, name string) ([]byte, error) {
	hash := sha256.New()
	buf := make([]byte, 1<<20) //nolint:mnd // 1 megabyte buffer

	for {
		index, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}

		if index == 0 {
			break
		}

		if _, err := flags.Write(buf[:index]); err != nil {
			return nil, err
		}

		if _, err := hash.Write(buf[:index]); err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	hashSum := hex.EncodeToString(hash.Sum(nil))
	if strings.ToLower(check) != hashSum {
		return nil, fmt.Errorf("hash mismatch for %s", name) //nolint:err113 // required name prevents static error
	}

	return data, nil
}

func extractMSI(flags *aflag.Flags, args ...string) error {
	if flags.MSIExtractVerbose {
		args = append(args, "-s")
		return process.Exec("./rust-msiexec.exe", args...)
	}

	return process.Exec("./rust-msiexec.exe", args...)
}
