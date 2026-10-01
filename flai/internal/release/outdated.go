package release

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// Outdated says that the running flai is older than the newest flai release
// tagged in the history of the project it serves (S-0181): it may not have
// rules or fields the project already uses.
type Outdated struct {
	Running string `json:"running"`
	Newest  string `json:"newest"`
	Upgrade string `json:"upgrade"`
	Message string `json:"message"`
}

// outdatedFor is how long an answer of FlaiOutdated is kept per project: the
// MCP inbox and the dashboard ask often, and a release lands rarely.
const outdatedFor = time.Minute

type outdatedAnswer struct {
	at  time.Time
	out *Outdated
}

var outdatedSeen sync.Map // root and running version → outdatedAnswer

// FlaiOutdated compares the running flai's version with the newest flai/vX.Y.Z
// tag reachable from the project's HEAD, and says so when it is older. nil
// when it is not, when running is not a release (a dev build), or when git
// cannot say. An answer is kept for a minute.
func FlaiOutdated(r execx.Runner, root, running string) *Outdated {
	if _, ok := buildinfo.Semver(running); !ok || r == nil {
		return nil
	}
	key := root + "\x00" + running
	if a, ok := outdatedSeen.Load(key); ok && time.Since(a.(outdatedAnswer).at) < outdatedFor {
		return a.(outdatedAnswer).out
	}
	out := flaiOutdated(r, root, running)
	outdatedSeen.Store(key, outdatedAnswer{at: time.Now(), out: out})
	return out
}

func flaiOutdated(r execx.Runner, root, running string) *Outdated {
	tags, err := r.Run(root, "git", "tag", "--merged", "HEAD", "--list", "flai/v*")
	if err != nil {
		return nil
	}
	newest := highestTag(strings.Split(strings.TrimSpace(tags), "\n"), "flai")
	if newest == (Version{}) || !buildinfo.Below(running, newest.String()) {
		return nil
	}
	return &Outdated{
		Running: strings.TrimPrefix(running, "v"),
		Newest:  newest.String(),
		Upgrade: buildinfo.UpgradeCommand,
		Message: fmt.Sprintf("this flai is %s, older than flai %s in the project's history, so it may lack rules and fields the project uses: upgrade it with %s",
			strings.TrimPrefix(running, "v"), newest, buildinfo.UpgradeCommand),
	}
}
