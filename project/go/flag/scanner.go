package flag

import (
	"bufio"
	"os"
	"strings"
)

func Parse(list []string, f *Flags) []string {
	replacements := map[string]string{
		"{HOST}":               f.HOST,
		"{TARGETX64}":          f.TARGETX64,
		"{TARGETX86}":          f.TARGETX86,
		"{TARGETARM}":          f.TARGETARM,
		"{TARGETARM_UPPER}":    strings.ToUpper(f.TARGETARM),
		"{TARGETARM64}":        f.TARGETARM64,
		"{TARGETARM64_UPPER}":  strings.ToUpper(f.TARGETARM64),
		"{MSVC_VERSION}":       f.MSVC_VERSION,
		"{MSVC_VERSION_LOCAL}": f.MSVC_VERSION_LOCAL,
		"{SDK_PID}":            f.SDK_PID,
	}

	parsed := []string{}
	for _, item := range list {
		for placeholder, value := range replacements {
			item = strings.ReplaceAll(item, placeholder, value)
		}
		parsed = append(parsed, item)
	}

	return parsed
}

func Scanner(file *os.File) ([]string, error) {
	entries := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			continue
		}
		entries = append(entries, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func IsFile(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return fileInfo.Mode().IsRegular(), nil
}
