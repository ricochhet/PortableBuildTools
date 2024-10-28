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
	"os"

	"github.com/ricochhet/minicommon/filesystem"
	"github.com/ricochhet/minicommon/logger"
	"github.com/ricochhet/minicommon/win32"
)

var (
	gitHash   string //nolint:gochecknoglobals // wontfix
	buildDate string //nolint:gochecknoglobals // wontfix
	buildOn   string //nolint:gochecknoglobals // wontfix
)

func printVersion() {
	logger.SharedLogger.Info(buildOn)
	logger.SharedLogger.Infof("GitHash: %s", gitHash)
	logger.SharedLogger.Infof("Build Date: %s", buildDate)
}

func logFile() *os.File { //nolint:mnd // wontfix
	file, err := os.OpenFile(filesystem.GetRelativePath("portablebuildtools_log.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal(err)
	}

	return file
}

//nolint:mnd // wontfix
func main() {
	logfile := logFile()
	defer func() {
		if err := logfile.Close(); err != nil {
			panic(err)
		}
	}()

	if err := win32.GuiConsoleHandle(os.Args, 1,
		func(_, cout, _ io.Writer) {
			logger.SharedLogger = logger.NewLogger(4, logger.InfoLevel, io.MultiWriter(logfile, cout), log.Lshortfile|log.LstdFlags)
			Cli(flags)
		},
		func(_, cout, _ io.Writer) {
			logger.SharedLogger = logger.NewLogger(4, logger.InfoLevel, io.MultiWriter(logfile, cout), log.Lshortfile|log.LstdFlags)
			logger.SharedLogger.Info("Initialized!")
			Gui(gitHash)
		}); err != nil {
		panic(err)
	}
}
