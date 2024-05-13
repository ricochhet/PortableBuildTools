package internal

import (
	"errors"

	aflag "github.com/ricochhet/portablebuildtools/flag"
)

var errMSIExtractMissing = errors.New("MSIExtract tool was not found")

func ExtractMSI(flags *aflag.Flags, args ...string) error {
	if exists, err := aflag.IsFile("./MSIExtract.exe"); err != nil || !exists {
		return errMSIExtractMissing
	}

	if flags.MSIExtractVerbose {
		args = append(args, "-s")
		return Exec("./MSIExtract.exe", args...)
	}

	return Exec("./MSIExtract.exe", args...)
}
