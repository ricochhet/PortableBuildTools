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
	"os"
	"syscall"

	"golang.org/x/sys/windows"
)

func AttachConsole() error {
	const ATTACH_PARENT_PROCESS = ^uintptr(0)
	proc := syscall.MustLoadDLL("kernel32.dll").MustFindProc("AttachConsole")
	r1, _, err := proc.Call(ATTACH_PARENT_PROCESS)
	if r1 == 0 {
		errno, ok := err.(syscall.Errno)
		if ok && errno == 6 {
			return nil
		}
		return err
	} else {
		return nil
	}
}

var wstdin, wstdout, wstderr *os.File

func AttachConsoleW() error {
	wstdin, wstdout, wstderr = os.Stdin, os.Stdout, os.Stderr
	if wstdin == nil || wstdout == nil || wstderr == nil {
		return fmt.Errorf("nil")
	}

	stdin, _ := syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
	stdout, _ := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	stderr, _ := syscall.GetStdHandle(syscall.STD_ERROR_HANDLE)

	var invalid syscall.Handle
	con := invalid

	if stdin == invalid || stdout == invalid || stderr == invalid {
		err := AttachConsole()
		if err != nil {
			return fmt.Errorf("attachconsole: %v", err)
		}

		if stdin == invalid {
			stdin, _ = syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
		}
		if stdout == invalid {
			stdout, _ = syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
			con = stdout
		}
		if stderr == invalid {
			stderr, _ = syscall.GetStdHandle(syscall.STD_ERROR_HANDLE)
			con = stderr
		}
	}

	if con != invalid {
		h := windows.Handle(con)
		var st uint32
		err := windows.GetConsoleMode(h, &st)
		if err != nil {
			return fmt.Errorf("GetConsoleMode: %v", err)
		}
		err = windows.SetConsoleMode(h, st&^windows.DISABLE_NEWLINE_AUTO_RETURN)
		if err != nil {
			return fmt.Errorf("SetConsoleMode: %v", err)
		}
	}

	if stdin != invalid {
		os.Stdin = os.NewFile(uintptr(stdin), "stdin")
	}
	if stdout != invalid {
		os.Stdout = os.NewFile(uintptr(stdout), "stdout")
	}
	if stderr != invalid {
		os.Stderr = os.NewFile(uintptr(stderr), "stderr")
	}
	return nil
}
