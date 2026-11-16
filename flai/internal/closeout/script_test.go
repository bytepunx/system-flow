package closeout

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// The close-out scripts the tests run: this repository's and the template's.
// The repository's calls scripts/flai.sh; the template's calls flai on PATH.
var closeOutScripts = []struct{ name, path string }{
	{"repository", "scripts/close-out.sh"},
	{"template", "template/root/scripts/close-out.sh"},
}

// stubFlai stands in for flai and scripts/flai.sh. It logs each call. The
// first verify --record-issues stops at sync, as flai verify does on a branch
// behind main; stream sync exits $SYNC_EXIT; every other call passes.
const stubFlai = `#!/bin/sh
echo "flai $*" >>"$STUB_LOG"
case "$1" in
  verify)
    case " $* " in
      *" --record-issues "*)
        n=$(cat "$STUB_COUNT" 2>/dev/null || echo 0)
        n=$((n + 1))
        echo "$n" >"$STUB_COUNT"
        if [ "$n" -eq 1 ]; then
          echo "verify: S-0001 stopped at sync (main has 3 commits the branch lacks)"
          exit 1
        fi
        ;;
    esac
    echo "verify: S-0001 passed"
    ;;
  stream) exit "${SYNC_EXIT:-0}" ;;
esac
`

// stubGit answers the checkout is on story/S-0001 with nothing uncommitted.
const stubGit = `#!/bin/sh
echo "git $*" >>"$STUB_LOG"
case "$1" in
  rev-parse) echo story/S-0001 ;;
esac
`

// stubEnv is scripts/env.sh reduced to what the close-out reads: ROOT.
const stubEnv = `ROOT="$(cd "$(dirname "$0")/.." && pwd)"
export ROOT
`

// runCloseOut copies the script at rel into a root of stubs and runs it for
// S-0001 with sh, answering the calls it logged and its last line.
func runCloseOut(t *testing.T, rel, syncExit string) (calls []string, last string) {
	t.Helper()
	sh, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("no sh on PATH")
	}
	script, err := os.ReadFile(filepath.Join("..", "..", "..", rel))
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	bin := filepath.Join(root, "stub-bin")
	for path, body := range map[string]string{
		"scripts/close-out.sh": string(script),
		"scripts/env.sh":       stubEnv,
		"scripts/flai.sh":      stubFlai,
		"stub-bin/flai":        stubFlai,
		"stub-bin/git":         stubGit,
	} {
		file := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file, []byte(body), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	log := filepath.Join(root, "calls.log")
	cmd := exec.Command(sh, filepath.Join(root, "scripts", "close-out.sh"), "S-0001", "-m", "close out")
	cmd.Dir = root
	cmd.Env = append(os.Environ(),
		"PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"),
		"STUB_LOG="+log,
		"STUB_COUNT="+filepath.Join(root, "verify.count"),
		"SYNC_EXIT="+syncExit,
	)
	out, _ := cmd.CombinedOutput()
	logged, err := os.ReadFile(log)
	if err != nil {
		t.Fatalf("the close-out called nothing: %v\n%s", err, out)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	return strings.Split(strings.TrimSpace(string(logged)), "\n"), lines[len(lines)-1]
}

// count is how many calls start with prefix.
func count(calls []string, prefix string) int {
	n := 0
	for _, c := range calls {
		if strings.HasPrefix(c, prefix) {
			n++
		}
	}
	return n
}

// A close-out asked for a branch behind main syncs and verifies again, rather
// than ending at its first flai verify's stop at sync (I-0383).
func TestACloseOutWhoseFirstVerifyStopsAtSyncSyncsAndVerifiesAgain(t *testing.T) {
	for _, s := range closeOutScripts {
		t.Run(s.name, func(t *testing.T) {
			calls, last := runCloseOut(t, s.path, "0")
			if n := count(calls, "flai stream sync S-0001"); n != 1 {
				t.Errorf("flai stream sync ran %d times, want 1; calls:\n%s", n, strings.Join(calls, "\n"))
			}
			if n := count(calls, "flai verify S-0001 --record-issues"); n != 2 {
				t.Errorf("flai verify --record-issues ran %d times, want 2; calls:\n%s", n, strings.Join(calls, "\n"))
			}
			if want := "close-out: S-0001 passed every step in round 2 of 3; ready to move to review"; last != want {
				t.Errorf("last line %q, want %q", last, want)
			}
		})
	}
}

// A close-out whose sync of a branch behind main fails stops at flai stream
// sync, in the round it synced for, and verifies nothing more.
func TestACloseOutWhoseSyncOfABranchBehindMainFailsStopsAtTheSync(t *testing.T) {
	for _, s := range closeOutScripts {
		t.Run(s.name, func(t *testing.T) {
			calls, last := runCloseOut(t, s.path, "1")
			if n := count(calls, "flai verify S-0001 --record-issues"); n != 1 {
				t.Errorf("flai verify --record-issues ran %d times, want 1; calls:\n%s", n, strings.Join(calls, "\n"))
			}
			if want := "close-out: S-0001 stopped at flai stream sync in round 2 of 3 (exit 1)"; !strings.HasPrefix(last, want) {
				t.Errorf("last line %q, want it to start %q", last, want)
			}
		})
	}
}
