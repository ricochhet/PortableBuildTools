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

	"github.com/ricochhet/minicommon/charmbracelet"
	"github.com/ricochhet/minicommon/filesystem"
	"github.com/ricochhet/minicommon/win32"
)

var (
	gitHash   string //nolint:gochecknoglobals // wontfix
	buildDate string //nolint:gochecknoglobals // wontfix
	buildOn   string //nolint:gochecknoglobals // wontfix
	debug     bool   //nolint:gochecknoglobals // wontfix
)

func printVersion() {
	charmbracelet.SharedLogger.Info(buildOn)
	charmbracelet.SharedLogger.Infof("GitHash: %s", gitHash)
	charmbracelet.SharedLogger.Infof("Build Date: %s", buildDate)
}

func logFile() *os.File { //nolint:mnd // wontfix
	file, err := os.OpenFile(filesystem.GetRelativePath("portablebuildtools_log.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal(err)
	}

	return file
}

func main() {
	logfile := logFile()
	defer func() {
		if err := logfile.Close(); err != nil {
			panic(err)
		}
	}()

	if err := win32.GuiConsoleHandle(os.Args, 1,
		func(_, cout, _ io.Writer) {
			Cli(flags, logfile, cout)
		},
		func(_, _, _ io.Writer) {
			Gui(gitHash, logfile)
		}, debug); err != nil {
		panic(err)
	}
}
