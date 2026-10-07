//go:build unix

package integration

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
)

// TestGolangciLintRunsBesideAnotherRunner holds the lock golangci-lint takes
// in the temp folder, as another story's run on the host would, and lints a
// small module with flai/.golangci.yaml. Without run.allow-parallel-runners
// golangci-lint waits 5s for the lock and exits 3 (I-0101, S-0308). It needs
// golangci-lint v2 on PATH, as scripts/integration.sh puts it there.
func TestGolangciLintRunsBesideAnotherRunner(t *testing.T) {
	if testing.Short() {
		t.Skip("runs golangci-lint")
	}
	lint, err := exec.LookPath("golangci-lint")
	if err != nil {
		t.Skip("golangci-lint not on PATH; run scripts/install-tools.sh")
	}
	config, err := filepath.Abs(filepath.Join("..", "..", ".golangci.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	tmp := t.TempDir()
	lock, err := os.Create(filepath.Join(tmp, "golangci-lint.lock"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = lock.Close() }()
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		t.Fatalf("holding the lock: %v", err)
	}

	mod := t.TempDir()
	write := func(name, text string) {
		if err := os.WriteFile(filepath.Join(mod, name), []byte(text), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	write("go.mod", "module example.com/lockcheck\n\ngo 1.26.0\n")
	write("lockcheck.go", "// Package lockcheck is linted while another runner holds the lock.\npackage lockcheck\n\n// One is one.\nfunc One() int { return 1 }\n")

	cmd := exec.Command(lint, "run", "--config", config, "./...")
	cmd.Dir = mod
	cmd.Env = append(os.Environ(), "TMPDIR="+tmp, "GOFLAGS=")
	out, err := cmd.CombinedOutput()
	var exit *exec.ExitError
	if errors.As(err, &exit) && exit.ExitCode() == 3 {
		t.Fatalf("golangci-lint gave up on the held lock (exit 3); flai/.golangci.yaml needs run.allow-parallel-runners:\n%s", out)
	}
	if err != nil {
		t.Fatalf("golangci-lint: %v\n%s", err, out)
	}
}
