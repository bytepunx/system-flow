//go:build windows

package components

import (
	"os/exec"
	"strconv"
	"syscall"
)

// viteBin is the name of the vite launcher a package manager installs in a
// sub-project's node_modules/.bin.
const viteBin = "vite.cmd"

// ownGroup starts cmd in a process group of its own, without a console.
func ownGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x00000200 | 0x08000000} // CREATE_NEW_PROCESS_GROUP | CREATE_NO_WINDOW
}

// terminateGroup stops the process and its tree: Windows has no gentler
// signal, so it is killGroup.
func terminateGroup(pid int) error { return killGroup(pid) }

// killGroup forces the process and everything it started to stop.
func killGroup(pid int) error {
	return exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(pid)).Run()
}
