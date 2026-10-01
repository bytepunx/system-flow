// Package buildinfo exposes version metadata injected at build time.
package buildinfo

import (
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
)

// Info is the resolved build metadata.
type Info struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
	Go      string `json:"go"`
}

// Get returns build metadata, filling commit and date from the Go build
// info when they were not injected by the linker.
func Get() Info {
	info := Info{Version: Version, Commit: Commit, Date: Date, Go: runtime.Version()}
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
