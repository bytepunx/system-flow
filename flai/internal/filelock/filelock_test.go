package filelock

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// The child in TestLockFreedWhenHolderDies is this test binary, told by
// childEnv which file to lock: it locks it, says so on stdout, and holds it
// until it is killed.
const childEnv = "FLAI_FILELOCK_TEST_HOLD"

func TestMain(m *testing.M) {
	path := os.Getenv(childEnv)
	if path == "" {
		os.Exit(m.Run())
	}
	if _, err := Lock(path, 0); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	fmt.Println("held")
	time.Sleep(time.Minute)
	os.Exit(0)
}

// TestLockSerializesWriters has many writers each read a counter, add one,
// and write it back under the lock: none loses another's increment.
func TestLockSerializesWriters(t *testing.T) {
	dir := t.TempDir()
	lock := filepath.Join(dir, "counter.lock")
	counter := filepath.Join(dir, "counter")
	if err := os.WriteFile(counter, []byte("0"), 0o600); err != nil {
		t.Fatal(err)
	}
	const writers, rounds = 16, 25
	var wg sync.WaitGroup
	errs := make(chan error, writers)
	for range writers {
		wg.Go(func() {
			for range rounds {
				if err := increment(lock, counter); err != nil {
					errs <- err
					return
				}
			}
		})
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatal(err)
	}
	if got := read(t, counter); got != writers*rounds {
		t.Fatalf("counter = %d, want %d: increments were lost", got, writers*rounds)
	}
}

func increment(lock, counter string) error {
	unlock, err := Lock(lock, 0)
	if err != nil {
		return err
	}
	defer unlock()
	b, err := os.ReadFile(counter)
	if err != nil {
		return err
	}
	n, err := strconv.Atoi(string(b))
	if err != nil {
		return fmt.Errorf("counter holds %q: %w", b, err)
	}
	return os.WriteFile(counter, []byte(strconv.Itoa(n+1)), 0o600)
}

func read(t *testing.T, path string) int {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	n, err := strconv.Atoi(string(b))
	if err != nil {
		t.Fatal(err)
	}
	return n
}

// TestLockTimesOutNamingFile: a lock still held when the wait passes is an
// error naming the file, and the lock is there to take once let go.
func TestLockTimesOutNamingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "folder", "x.lock")
	unlock, err := Lock(path, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	if _, err := Lock(path, 50*time.Millisecond); err == nil {
		t.Fatal("a second Lock took a held lock")
	} else if !strings.Contains(err.Error(), path) {
		t.Fatalf("error %q does not name %s", err, path)
	}
	if waited := time.Since(start); waited < 50*time.Millisecond {
		t.Fatalf("gave up after %s, before its wait of 50ms", waited)
	}
	unlock()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("unlock removed the lock file: %v", err)
	}
	again, err := Lock(path, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("lock not free after unlock: %v", err)
	}
	again()
}

// TestLockFreedWhenHolderDies kills a child process that holds the lock;
// the lock is free at once, with no wait for its age.
func TestLockFreedWhenHolderDies(t *testing.T) {
	path := filepath.Join(t.TempDir(), "held.lock")
	cmd := exec.Command(os.Args[0], "-test.run=^$")
	cmd.Env = append(os.Environ(), childEnv+"="+path)
	cmd.Stderr = os.Stderr
	out, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill(); _ = cmd.Wait() })
	line, err := bufio.NewReader(out).ReadString('\n')
	if err != nil || strings.TrimSpace(line) != "held" {
		t.Fatalf("child did not take the lock: %q, %v", line, err)
	}

	if _, err := Lock(path, 50*time.Millisecond); err == nil {
		t.Fatal("took the lock while the child held it")
	} else if !strings.Contains(err.Error(), path) {
		t.Fatalf("error %q does not name %s", err, path)
	}

	if err := cmd.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	_ = cmd.Wait()
	start := time.Now()
	unlock, err := Lock(path, 500*time.Millisecond)
	if err != nil {
		t.Fatalf("lock not freed by the holder's death: %v", err)
	}
	unlock()
	t.Logf("took the dead holder's lock after %s", time.Since(start))
}
