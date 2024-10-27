package main

import (
	"github.com/ricochhet/minicommon/logger"
	aflag "github.com/ricochhet/portablebuildtools/flag"
)

func Cli(flags *aflag.Flags) {
	errCh := make(chan error, 3)

	if flags.Version {
		printVersion()
		return
	}

	writeEnvironmentsWerr(errCh)
	for err := range errCh {
		if err != nil {
			logger.SharedLogger.Fatalf("FATAL: %v", err)
		}
	}

	runWerr(errCh)
	for err := range errCh {
		if err != nil {
			logger.SharedLogger.Fatalf("FATAL: %v", err)
		}
	}
}
