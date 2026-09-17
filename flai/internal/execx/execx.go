// Package execx runs external programs (git, docker) behind an interface so
// commands can be tested without them. See ADR-0010.
package execx

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// Runner executes a program and returns its combined output.
type Runner interface {
	Run(dir, name string, args ...string) (string, error)
	// RunInput is Run with data on standard input, for secrets that must
	// not appear in arguments (docker login --password-stdin).
	RunInput(dir, name, input string, args ...string) (string, error)
	LookPath(name string) (string, error)
}

// System is the real Runner.
type System struct{}

// Run executes name with args in dir and returns combined stdout and stderr.
func (s System) Run(dir, name string, args ...string) (string, error) {
	return s.RunInput(dir, name, "", args...)
}

// RunInput executes name with args in dir, feeding input on stdin, and
// returns combined stdout and stderr.
func (System) RunInput(dir, name, input string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	if input != "" {
		cmd.Stdin = strings.NewReader(input)
	}
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	out := strings.TrimSpace(buf.String())
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return out, fmt.Errorf("%s %s: %w\n%s", name, strings.Join(args, " "), err, out)
		}
		return out, fmt.Errorf("%s: %w", name, err)
	}
	return out, nil
}

// LookPath reports where name is on PATH.
func (System) LookPath(name string) (string, error) {
	return exec.LookPath(name)
}

// Require returns a helpful error when name is not installed.
func Require(r Runner, name, hint string) error {
	if _, err := r.LookPath(name); err != nil {
		return fmt.Errorf("%s is required but not found on PATH. %s", name, hint)
	}
	return nil
}
