/*
 * PortableBuildTools
 * Copyright (C) 2024 PortableBuildTools contributors
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"io"
	"log"

	"github.com/ricochhet/minicommon/logger"
	aflag "github.com/ricochhet/portablebuildtools/flag"
)

//nolint:mnd // wontfix
func Cli(flags *aflag.Flags, logfile, cout io.Writer) {
	logger.SharedLogger = logger.NewLogger(4, logger.InfoLevel, io.MultiWriter(logfile, cout), log.Lshortfile|log.LstdFlags)

	errCh := make(chan error, 3) //nolint:mnd // wontfix

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
