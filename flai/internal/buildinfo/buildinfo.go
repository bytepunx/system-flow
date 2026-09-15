// Package buildinfo exposes version metadata injected at build time.
package buildinfo

import (
	"runtime"
	"runtime/debug"
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
