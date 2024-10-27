//go:build windows

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

package win32

import (
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/ricochhet/portablebuildtools/thirdparty/ansi"
)

func AllocConsole() (aIn, aOut, aErr io.Writer, e error) {
	kernal23 := syscall.NewLazyDLL("kernel32.dll")
	allocConsole := kernal23.NewProc("AllocConsole")
	r0, _, err0 := syscall.SyscallN(allocConsole.Addr(), 0, 0, 0, 0)
	if r0 == 0 {
		return nil, nil, nil, fmt.Errorf("could not allocate console: %s. check build flags", err0)
	}

	hIn, err1 := syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
	hOut, err2 := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	hErr, err3 := syscall.GetStdHandle(syscall.STD_ERROR_HANDLE)

	if err1 != nil {
		return nil, nil, nil, err1
	}

	if err2 != nil {
		return nil, nil, nil, err2
	}

	if err3 != nil {
		return nil, nil, nil, err3
	}

	stdinFile := os.NewFile(uintptr(hIn), "/dev/stdin")
	stdoutFile := os.NewFile(uintptr(hOut), "/dev/stdout")
	stderrFile := os.NewFile(uintptr(hErr), "/dev/stderr")

	aStdin := ansi.NewAnsiStdoutW(stdinFile)
	aStdout := ansi.NewAnsiStdoutW(stdoutFile)
	aStderr := ansi.NewAnsiStdoutW(stderrFile)

	os.Stdin = stdinFile
	os.Stdout = stdoutFile
	os.Stderr = stderrFile

	return aStdin, aStdout, aStderr, nil
}
