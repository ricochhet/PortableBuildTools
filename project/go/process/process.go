package process

import (
	"os"
	"os/exec"
)

func Exec(command string, args ...string) error {
	cmd := exec.Command(command, args...)

	// Attach stdout and stderr for output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
