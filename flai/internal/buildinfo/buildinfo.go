// Package buildinfo exposes version metadata injected at build time.
package buildinfo

import (
	_ "embed"
	"fmt"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
)

// Set via -ldflags "-X github.com/bytepunx/system-flow/flai/internal/buildinfo.Version=...".
var (
	Version = "dev"
	Commit  = ""
	Date    = ""
	// ReleaseStatement is the release stamp's statement, "flai <version>
	// <commit>" on one line, which the release workflow signs (S-0235).
	ReleaseStatement = ""
	// ReleaseSignature is the statement's signature by the release key,
	// base64 of its ASN.1 DER form, as cosign sign-blob writes it.
	ReleaseSignature = ""
)

// Info is the resolved build metadata.
type Info struct {
	Version       string `json:"version"`
	Commit        string `json:"commit"`
	Date          string `json:"date"`
	Go            string `json:"go"`
	TaskDoneRules int    `json:"task_done_rules"`
}

// The task-done rules number (ADR-0180): task-done-rules.txt holds one
// positive integer, raised by one in any change to what flai task done
// commits, messages, syncs, moves, logs, adds to the touches, or checks. A
// flai/vX.Y.Z tag records it as that file at the commit it tags, so that
// flai task done refuses an older flai only when a release raised it.
//
//go:embed task-done-rules.txt
var taskDoneRulesText string

// TaskDoneRulesPath is where the task-done rules number is kept, from the
// root of the repository that tags flai releases (ADR-0180).
const TaskDoneRulesPath = "flai/internal/buildinfo/task-done-rules.txt"

// TaskDoneRules is this flai's task-done rules number (ADR-0180).
var TaskDoneRules = mustTaskDoneRules(taskDoneRulesText)

// ParseTaskDoneRules reads a task-done rules number from the text of
// task-done-rules.txt: one positive integer, with surrounding space allowed.
func ParseTaskDoneRules(text string) (int, bool) {
	n, err := strconv.Atoi(strings.TrimSpace(text))
	if err != nil || n < 1 {
		return 0, false
	}
	return n, true
}

func mustTaskDoneRules(text string) int {
	n, ok := ParseTaskDoneRules(text)
	if !ok {
		panic(fmt.Sprintf("buildinfo: %s holds %q, not a positive integer: write the task-done rules number there (ADR-0180)", TaskDoneRulesPath, text))
	}
	return n
}

// Get returns build metadata, filling commit and date from the Go build
// info when they were not injected by the linker.
func Get() Info {
	info := Info{Version: Version, Commit: Commit, Date: Date, Go: runtime.Version(), TaskDoneRules: TaskDoneRules}
	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				if info.Commit == "" {
					info.Commit = s.Value
				}
			case "vcs.time":
				if info.Date == "" {
					info.Date = s.Value
				}
			}
		}
	}
	if info.Commit == "" {
		info.Commit = "unknown"
	}
	if info.Date == "" {
		info.Date = "unknown"
	}
	return info
}

// UpgradeCommand is what brings the host's flai up to date.
const UpgradeCommand = "flai host upgrade (flai self-upgrade where no flai host runs)"

// Semver returns a release version's major, minor, and patch numbers from
// X.Y.Z, with or without a leading v; a pre-release or build suffix is left
// out. Anything else, such as dev, is not a release.
func Semver(v string) ([3]int, bool) {
	var out [3]int
	core, _, _ := strings.Cut(strings.TrimPrefix(strings.TrimSpace(v), "v"), "-")
	core, _, _ = strings.Cut(core, "+")
	parts := strings.Split(core, ".")
	if len(parts) != 3 {
		return out, false
	}
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return out, false
		}
		out[i] = n
	}
	return out, true
}

// Bare reports whether v is a release written as a bare X.Y.Z: no leading
// v, no suffix, and no leading zeros (ADR-0117).
func Bare(v string) bool {
	n, ok := Semver(v)
	return ok && fmt.Sprintf("%d.%d.%d", n[0], n[1], n[2]) == v
}

// Below reports whether version is a release older than other. A version
// that is not a release, such as a dev build, is never below anything, and
// nothing is below a version that is not one.
func Below(version, other string) bool {
	a, ok := Semver(version)
	b, ok2 := Semver(other)
	if !ok || !ok2 {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return false
}
