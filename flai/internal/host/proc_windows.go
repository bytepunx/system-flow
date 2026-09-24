//go:build windows

package host

import (
	"os"
	"os/exec"
)

// alive cannot probe with a signal on Windows; the fresh state file is the
// evidence there.
func alive(pid int) bool { return pid > 0 }

// terminate stops a child; Windows has no gentler signal to send.
func terminate(p *os.Process) error { return p.Kill() }

// Reexec starts the flai at exe with the same arguments and ends this
// process: Windows cannot replace a running image.
func Reexec(exe string) error {
	cmd := exec.Command(exe, os.Args[1:]...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	os.Exit(0)
	return nil
}
