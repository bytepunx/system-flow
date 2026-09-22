//go:build windows

package serve

import (
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

// Alive cannot probe with a signal on Windows; the fresh status
// file is the evidence there.
func Alive(pid int) bool { return pid > 0 }

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

// TerminateGroup stops the process and its tree (S-0082): Windows has no
// gentler signal, so this is the same as KillGroup.
func TerminateGroup(pid int) error { return KillGroup(pid) }

// KillGroup forces the process and everything it spawned to stop.
func KillGroup(pid int) error {
	return exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(pid)).Run()
}
