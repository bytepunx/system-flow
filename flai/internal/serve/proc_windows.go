//go:build windows

package serve

import (
	"os"
	"os/exec"
	"syscall"
)

// processAlive cannot probe with a signal on Windows; the fresh status
// file is the evidence there.
func processAlive(pid int) bool { return pid > 0 }

// Detach starts cmd in its own process group, without a console.
func Detach(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x00000200 | 0x08000000} // CREATE_NEW_PROCESS_GROUP | CREATE_NO_WINDOW
}

// Terminate stops the process; Windows has no gentler signal to send.
func Terminate(pid int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Kill()
}
