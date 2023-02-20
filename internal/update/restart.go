package update

import (
	"os"
	"os/exec"
)

func RestartSelf() (err error) {
	var cmd *exec.Cmd
	if len(os.Args) == 1 {
		cmd = exec.Command(os.Args[0])
	} else {
		cmd = exec.Command(os.Args[0], os.Args[1:]...)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Start(); err != nil {
		return
	}
	err = cmd.Process.Release()
	return
}
