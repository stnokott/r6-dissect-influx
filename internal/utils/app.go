// go:build windows

package utils

import (
	"os"
	"os/exec"
)

func RestartApp() error {
	self, err := os.Executable()
	if err != nil {
		return err
	}
	args := os.Args
	env := os.Environ()

	// on Linux, we would use syscall.Exec
	cmd := exec.Command(self, args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = env
	if err := cmd.Start(); err == nil {
		os.Exit(0)
	} else {
		return err
	}
	return nil
}
