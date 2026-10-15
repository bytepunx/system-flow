package serve

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Catching up with the projects it serves (S-0387, ADR-0158, I-0141). Every
// minute, and as soon as it serves a project it had not, flai serve asks
// release.FlaiOutdated whether its flai is older than the newest flai/vX.Y.Z
// tag in each served project's history. For a project behind it with the
// shell-only host action auto-upgrade on, it asks the host it runs under to
// install that release, the highest where several are behind, and restart
// on it; never while a dashboard write runs in its process, which the
// restart would cut, and for a release the host refused, no sooner than 15
// minutes after the last ask. A project behind it with auto-upgrade off, or
// with it on while flai serve runs under no host, or whose release is still
// not installed an hour after the first ask, gets one thread on its
// system-flow.yaml, once per release, recorded under the serve directory so
// that no later look and no restart opens it again.

const (
	// upgradeEvery is how often the served projects are looked at again.
	upgradeEvery = time.Minute
	// upgradeRetry is the least time between two asks for one release.
	upgradeRetry = 15 * time.Minute
	// upgradeTellAfter is how long a release asked for may stay not
	// installed before the operator is told on a thread.
	upgradeTellAfter = time.Hour
	// upgradeAuthor opens the thread the operator is asked on.
	upgradeAuthor = "flai serve"
	// upgradeTimeout bounds one ask of the host: it downloads and verifies
	// the release before it answers.
	upgradeTimeout = 6 * time.Minute
)

// HostUpgrade asks the host flai serve runs under to install the flai
// release version, as flai host upgrade --version does, and answers whether
// it installed it and so restarts.
type HostUpgrade func(ctx context.Context, version string) (restarting bool, err error)

func (d Dir) upgradeAsksFile() string { return filepath.Join(string(d), "upgrade-asks.json") }

// upgradeAsks reads the newest flai release the operator was asked to
// upgrade to, by project root; none is an empty map.
func (d Dir) upgradeAsks() (map[string]string, error) {
	out := map[string]string{}
	data, err := os.ReadFile(d.upgradeAsksFile())
	if errors.Is(err, os.ErrNotExist) {
		return out, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("%s: %w", d.upgradeAsksFile(), err)
	}
	return out, nil
}

// recordUpgradeAsk records that the operator was asked to upgrade to version
// for the project at root.
func (d Dir) recordUpgradeAsk(root, version string) error {
	all, err := d.upgradeAsks()
	if err != nil {
		return err
	}
	all[root] = version
	return d.write(d.upgradeAsksFile(), all)
}

// upgradeAttempt is what became of the asks for one release.
type upgradeAttempt struct {
	first, last time.Time
	refusal     string // the host's last refusal, "" when it took the ask
}

// upgrader looks at the served projects' flai releases and catches up or
// asks; its dependencies are funcs so that it is tested without a host, git,
// or a repository.
type upgrader struct {
	now      func() time.Time
	log      *slog.Logger
	dir      Dir
	outdated func(root string) *release.Outdated
	auto     func(root string) bool
	upgrade  HostUpgrade // nil when flai serve runs under no host
	writing  func() int
	tell     func(e Entry, title, text string) (thread string, err error)
	journal  func(hostapi.Entry)

	busy     atomic.Bool
	wg       sync.WaitGroup
	last     time.Time
	seen     map[string]bool
	attempts map[string]*upgradeAttempt // by release
	told     map[string]string          // root → release asked about, as recorded
}

func newUpgrader(o Options) *upgrader {
	git := gitOf(o)
	return &upgrader{
		now:      o.Now,
		log:      o.Logger,
		dir:      o.Dir,
		outdated: func(root string) *release.Outdated { return release.FlaiOutdated(git, root, o.Version) },
		auto: func(root string) bool {
			return o.Host.Enabled != nil && o.Host.Enabled(hostapi.ActionAutoUpgrade, root)
		},
		upgrade: o.Upgrade,
		writing: hostapi.WritesRunning,
		tell:    upgradeThread(o.Now),
		journal: o.Host.Record,
	}
}

// upgradeThread opens the thread that asks the operator to upgrade, on the
// project's manifest, as flai serve.
func upgradeThread(now func() time.Time) func(e Entry, title, text string) (string, error) {
	return func(e Entry, title, text string) (string, error) {
		repo, err := workitem.Open(e.Root)
		if err != nil {
			return "", err
		}
		th, err := threads.New(repo, threads.NewOptions{Title: title, On: manifest.File, Author: upgradeAuthor, Text: text, Now: now()})
		if err != nil {
			return "", err
		}
		return th.ID, nil
	}
}

// soon looks in the background unless a look still runs: an ask waits on
// the host for minutes, and the serve loop that calls it must not.
func (u *upgrader) soon(ctx context.Context, served []Entry) {
	if !u.busy.CompareAndSwap(false, true) {
		return
	}
	u.wg.Add(1)
	go func() {
		defer u.wg.Done()
		defer u.busy.Store(false)
		u.look(ctx, served)
	}()
}

// wait waits for a look that runs.
func (u *upgrader) wait() { u.wg.Wait() }

// due says whether the projects are looked at now: a minute after the last
// look, or at once for a project not looked at before.
func (u *upgrader) due(now time.Time, served []Entry) bool {
	if u.seen == nil {
		u.seen = map[string]bool{}
	}
	due := u.last.IsZero() || now.Sub(u.last) >= upgradeEvery
	for _, e := range served {
		if !u.seen[e.Root] {
			u.seen[e.Root] = true
			due = true
		}
	}
	if due {
		u.last = now
	}
	return due
}

type behindProject struct {
	e   Entry
	out *release.Outdated
}

// look catches up with, or asks about, every served project whose history
// holds a newer flai release than the one serving it, when it is due.
func (u *upgrader) look(ctx context.Context, served []Entry) {
	now := u.now()
	if !u.due(now, served) {
		return
	}
	var auto []behindProject
	best := -1 // the one behind the highest release, of auto
	for _, e := range served {
		out := u.outdated(e.Root)
		if out == nil {
			continue
		}
		switch {
		case !u.auto(e.Root):
			u.ask(e, out, "")
		case u.upgrade == nil:
			u.ask(e, out, "auto-upgrade is on for this project, but flai serve runs outside flai host, so there is no host for it to ask")
		default:
			auto = append(auto, behindProject{e, out})
			if best < 0 || buildinfo.Below(auto[best].out.Newest, out.Newest) {
				best = len(auto) - 1
			}
		}
	}
	if best < 0 {
		return
	}
	target := auto[best]
	version := target.out.Newest
	att := u.attempts[version]
	switch {
	case att != nil && now.Sub(att.last) < upgradeRetry:
		// asked lately: the host restarts when it installs, or refused
	case u.writing() > 0:
		u.log.Info("flai upgrade waits for a dashboard write", "component", "serve", "running", target.out.Running, "newest", version, "writes", u.writing())
	default:
		att = u.request(ctx, now, target)
	}
	if att == nil || now.Sub(att.first) < upgradeTellAfter {
		return
	}
	why := fmt.Sprintf("auto-upgrade is on for this project, and flai serve first asked flai host for flai %s at %s, but it is still not installed", version, att.first.UTC().Format(time.RFC3339))
	if att.refusal != "" {
		why += "; the host's last answer: " + att.refusal
	}
	for _, b := range auto {
		u.ask(b.e, b.out, why)
	}
}

// request asks the host to install b's newest release, logs and journals
// what it answered, and keeps the ask; nil when serve is stopping and so
// nothing was asked.
func (u *upgrader) request(ctx context.Context, now time.Time, b behindProject) *upgradeAttempt {
	version := b.out.Newest
	actx, cancel := context.WithTimeout(ctx, upgradeTimeout)
	defer cancel()
	restarting, err := u.upgrade(actx, version)
	if ctx.Err() != nil {
		return u.attempts[version]
	}
	if err == nil && !restarting {
		err = fmt.Errorf("the host installed nothing for flai %s", version)
	}
	if u.attempts == nil {
		u.attempts = map[string]*upgradeAttempt{}
	}
	att := u.attempts[version]
	if att == nil {
		att = &upgradeAttempt{first: now}
		u.attempts[version] = att
	}
	att.last = now
	entry := hostapi.Entry{At: now.UTC().Format(time.RFC3339), Action: hostapi.ActionAutoUpgrade, Method: "serve.upgrade", Project: b.e.Key, Root: b.e.Root, By: upgradeAuthor}
	if err != nil {
		att.refusal = err.Error()
		u.log.Warn("flai upgrade refused", "component", "serve", "project", b.e.Key, "running", b.out.Running, "newest", version, "err", att.refusal)
		entry.Outcome, entry.Detail = "failed", fmt.Sprintf("asked flai host to upgrade flai %s to %s: %s", b.out.Running, version, att.refusal)
	} else {
		att.refusal = ""
		u.log.Info("flai upgrade asked", "component", "serve", "project", b.e.Key, "running", b.out.Running, "newest", version)
		entry.Outcome, entry.Detail = "done", fmt.Sprintf("asked flai host to upgrade flai %s to %s; it installed it and restarts", b.out.Running, version)
	}
	if u.journal != nil {
		u.journal(entry)
	}
	return att
}

// ask opens the thread that asks the operator to upgrade flai for the
// project e to out's newest release, unless one was opened for that release
// or a newer one; why says why flai serve did not upgrade it itself, "" when
// auto-upgrade is off.
func (u *upgrader) ask(e Entry, out *release.Outdated, why string) {
	if u.told == nil {
		u.told = map[string]string{}
	}
	was := u.told[e.Root]
	if was == "" {
		asked, err := u.dir.upgradeAsks()
		if err != nil {
			u.log.Warn("flai upgrade asks unreadable", "component", "serve", "project", e.Key, "err", err.Error())
			return
		}
		was = asked[e.Root]
	}
	if was != "" && !buildinfo.Below(was, out.Newest) {
		u.told[e.Root] = was
		return
	}
	id, err := u.tell(e, upgradeTitle(out), upgradeText(out, why))
	if err != nil {
		u.log.Warn("flai upgrade not asked on a thread", "component", "serve", "project", e.Key, "running", out.Running, "newest", out.Newest, "err", err.Error())
		return
	}
	u.told[e.Root] = out.Newest
	if err := u.dir.recordUpgradeAsk(e.Root, out.Newest); err != nil {
		u.log.Warn("flai upgrade ask not recorded", "component", "serve", "project", e.Key, "newest", out.Newest, "file", u.dir.upgradeAsksFile(), "err", err.Error())
	}
	u.log.Info("flai upgrade asked on a thread", "component", "serve", "project", e.Key, "running", out.Running, "newest", out.Newest, "thread", id)
}

func upgradeTitle(out *release.Outdated) string {
	return fmt.Sprintf("The host's flai %s is older than flai %s in this project's history", out.Running, out.Newest)
}

// upgradeText asks the operator to upgrade: with auto-upgrade, why says why
// flai serve did not; without, it offers auto-upgrade.
func upgradeText(out *release.Outdated, why string) string {
	text := fmt.Sprintf("flai serve, and the MCP server it has flai host keep for this project, run flai %[1]s, older than flai %[2]s, the newest release tagged in this project's history. They may lack rules and fields the project already uses, and flai task done refuses while the minor version differs (ADR-0142).\n\n", out.Running, out.Newest)
	if why != "" {
		text += why + ".\n\n"
	}
	text += fmt.Sprintf("Upgrade it in a shell on the host with: flai host upgrade --version %s", out.Newest)
	if why == "" {
		text += "\n\nOr have flai serve upgrade the host whenever this project's history gets ahead of it: run flai serve enable auto-upgrade in the project (ADR-0158)"
	}
	return text + ".\n\nflai serve asks once for each release: a newer one asks again."
}
