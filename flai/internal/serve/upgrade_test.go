package serve

import (
	"context"
	"errors"
	"log/slog"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/release"
)

// upgradeLab is flai serve's upgrader with every dependency in hand: a
// clock, each project's newest release in its history, its auto-upgrade
// action, the host, the writes in flight, and the threads opened.
type upgradeLab struct {
	t        *testing.T
	dir      Dir
	at       time.Time
	running  string
	newest   map[string]string // root → newest release in its history
	auto     map[string]bool
	refuse   error // what the host answers an ask with
	writes   int
	asked    []string // the releases the host was asked for
	threads  []string // root: title, then text, of each thread opened
	journal  []hostapi.Entry
	logs     syncBuffer
	underNo  bool // flai serve runs under no host
	upgrader *upgrader
}

func newUpgradeLab(t *testing.T) *upgradeLab {
	lab := &upgradeLab{t: t, dir: DirFor(filepath.Join(t.TempDir(), "config.json")), at: time.Date(2026, 10, 15, 9, 0, 0, 0, time.UTC),
		running: "1.48.0", newest: map[string]string{}, auto: map[string]bool{}}
	lab.restart()
	return lab
}

// restart is a new flai serve on the same serve directory: it keeps nothing
// in memory, only what it recorded.
func (lab *upgradeLab) restart() {
	u := &upgrader{now: func() time.Time { return lab.at }, log: slog.New(slog.NewTextHandler(&lab.logs, nil)), dir: lab.dir,
		outdated: func(root string) *release.Outdated {
			if n := lab.newest[root]; n != "" && n != lab.running {
				return &release.Outdated{Running: lab.running, Newest: n}
			}
			return nil
		},
		auto:    func(root string) bool { return lab.auto[root] },
		writing: func() int { return lab.writes },
		tell: func(e Entry, title, text string) (string, error) {
			lab.threads = append(lab.threads, e.Root+": "+title+"\n"+text)
			return "TH-0001", nil
		},
		journal: func(e hostapi.Entry) { lab.journal = append(lab.journal, e) },
	}
	if !lab.underNo {
		u.upgrade = func(_ context.Context, version string) (bool, error) {
			lab.asked = append(lab.asked, version)
			return lab.refuse == nil, lab.refuse
		}
	}
	lab.upgrader = u
}

func (lab *upgradeLab) look(roots ...string) {
	served := make([]Entry, 0, len(roots))
	for _, r := range roots {
		served = append(served, Entry{Key: r, Root: r})
	}
	lab.upgrader.look(context.Background(), served)
}

func (lab *upgradeLab) after(d time.Duration) { lab.at = lab.at.Add(d) }

// I-0141, S-0387: a flai serve already running when the project's history
// gains a newer flai release finds it on its next minute's look and, with
// auto-upgrade on, asks the host for that release once, logs it, and
// journals it as a host action.
func TestARunningServeUpgradesToAReleaseTaggedAfterItStarted(t *testing.T) {
	lab := newUpgradeLab(t)
	lab.auto["harbour"] = true
	lab.newest["harbour"] = "1.48.0"
	lab.look("harbour")
	if len(lab.asked) != 0 || len(lab.threads) != 0 {
		t.Fatalf("a serve as new as its project asked %v and told %v", lab.asked, lab.threads)
	}
	lab.newest["harbour"] = "1.49.0" // the release lands while it runs
	lab.after(30 * time.Second)
	lab.look("harbour")
	if len(lab.asked) != 0 {
		t.Fatalf("looked again before the minute was up: %v", lab.asked)
	}
	lab.after(30 * time.Second)
	lab.look("harbour")
	if strings.Join(lab.asked, " ") != "1.49.0" {
		t.Fatalf("the host was asked for %v, want 1.49.0", lab.asked)
	}
	logs := lab.logs.String()
	if !strings.Contains(logs, `level=INFO msg="flai upgrade asked"`) || !strings.Contains(logs, "running=1.48.0") || !strings.Contains(logs, "newest=1.49.0") {
		t.Errorf("the ask is logged with both versions:\n%s", logs)
	}
	if len(lab.journal) != 1 || lab.journal[0].Action != hostapi.ActionAutoUpgrade || lab.journal[0].Outcome != "done" || lab.journal[0].Root != "harbour" {
		t.Errorf("the ask is journalled as the auto-upgrade host action: %+v", lab.journal)
	}
	if len(lab.threads) != 0 {
		t.Errorf("an upgrade asked for also asked the operator: %v", lab.threads)
	}
}

// ADR-0158: a project flai serve starts to serve is looked at at once, not
// a minute later.
func TestANewlyServedProjectIsLookedAtAtOnce(t *testing.T) {
	lab := newUpgradeLab(t)
	lab.auto["quay"] = true
	lab.look("harbour")
	lab.newest["quay"] = "1.49.0"
	lab.after(time.Second)
	lab.look("harbour", "quay")
	if strings.Join(lab.asked, " ") != "1.49.0" {
		t.Errorf("the host was asked for %v when quay was first served", lab.asked)
	}
}

// ADR-0158: where several projects with auto-upgrade on are behind, the host
// is asked once, for the highest of their releases; one with auto-upgrade
// off is asked about on its own thread.
func TestTheHighestReleaseIsAskedFor(t *testing.T) {
	lab := newUpgradeLab(t)
	lab.auto["harbour"], lab.auto["quay"] = true, true
	lab.newest["harbour"], lab.newest["quay"], lab.newest["dock"] = "1.49.2", "1.50.0", "1.51.0"
	lab.look("harbour", "quay", "dock")
	if strings.Join(lab.asked, " ") != "1.50.0" {
		t.Errorf("the host was asked for %v, want the highest of the auto-upgrade projects', 1.50.0", lab.asked)
	}
	if len(lab.threads) != 1 || !strings.HasPrefix(lab.threads[0], "dock: ") {
		t.Errorf("the project with auto-upgrade off is asked about on a thread: %v", lab.threads)
	}
}

// ADR-0158: no upgrade is asked while a dashboard write runs in flai serve's
// process; the next minute's look asks once it has ended.
func TestTheUpgradeWaitsForADashboardWrite(t *testing.T) {
	lab := newUpgradeLab(t)
	lab.auto["harbour"] = true
	lab.newest["harbour"] = "1.49.0"
	lab.writes = 1
	lab.look("harbour")
	if len(lab.asked) != 0 {
		t.Fatalf("asked for %v while a write ran", lab.asked)
	}
	if !strings.Contains(lab.logs.String(), `msg="flai upgrade waits for a dashboard write"`) {
		t.Errorf("the wait is logged:\n%s", lab.logs.String())
	}
	lab.writes = 0
	lab.after(time.Minute)
	lab.look("harbour")
	if strings.Join(lab.asked, " ") != "1.49.0" {
		t.Errorf("once the write ended the host was asked for %v", lab.asked)
	}
}

// ADR-0158: a refused release, tagged but not yet published say, is logged
// at warn and asked for again no sooner than 15 minutes later; still not
// installed an hour after the first ask, the operator is asked on a thread
// naming the last refusal, once, and the asks go on every 15 minutes.
func TestARefusedUpgradeIsAskedAgainEveryQuarterHourAndToldAfterAnHour(t *testing.T) {
	lab := newUpgradeLab(t)
	lab.auto["harbour"] = true
	lab.newest["harbour"] = "1.49.0"
	lab.refuse = errors.New("flai 1.49.0 is not published; published are 1.48.0")
	lab.look("harbour")
	logs := lab.logs.String()
	if len(lab.asked) != 1 || !strings.Contains(logs, `level=WARN msg="flai upgrade refused"`) || !strings.Contains(logs, "not published") {
		t.Fatalf("the refusal is logged at warn: asked %v\n%s", lab.asked, logs)
	}
	if len(lab.journal) != 1 || lab.journal[0].Outcome != "failed" {
		t.Errorf("the refusal is journalled: %+v", lab.journal)
	}
	for range 14 {
		lab.after(time.Minute)
		lab.look("harbour")
	}
	if len(lab.asked) != 1 {
		t.Fatalf("asked again within 15 minutes: %v", lab.asked)
	}
	lab.after(time.Minute)
	lab.look("harbour")
	if len(lab.asked) != 2 {
		t.Fatalf("not asked again 15 minutes later: %v", lab.asked)
	}
	for range 44 {
		lab.after(time.Minute)
		lab.look("harbour")
	}
	if len(lab.threads) != 0 {
		t.Fatalf("told before the hour was up: %v", lab.threads)
	}
	lab.refuse = errors.New("flai 1.49.0 has no signature that verifies")
	lab.after(time.Minute)
	lab.look("harbour")
	if len(lab.asked) != 5 {
		t.Errorf("asked %d times in the hour, want at the start and every 15 minutes: %v", len(lab.asked), lab.asked)
	}
	if len(lab.threads) != 1 || !strings.Contains(lab.threads[0], "no signature that verifies") || !strings.Contains(lab.threads[0], "flai host upgrade --version 1.49.0") {
		t.Fatalf("an hour on, the operator is told the last refusal: %v", lab.threads)
	}
	for range 30 {
		lab.after(time.Minute)
		lab.look("harbour")
	}
	if len(lab.threads) != 1 || len(lab.asked) != 7 {
		t.Errorf("the thread is opened once and the asks go on every 15 minutes: %d threads, asked %v", len(lab.threads), lab.asked)
	}
}

// ADR-0158: with auto-upgrade off, the operator is asked on one thread per
// project and release, naming both versions and both ways to catch up; no
// later look opens it again, nor does a restart of flai serve, which reads
// the record; a newer release opens a new one.
func TestWithAutoUpgradeOffTheOperatorIsAskedOncePerRelease(t *testing.T) {
	lab := newUpgradeLab(t)
	lab.newest["harbour"] = "1.49.0"
	lab.look("harbour")
	if len(lab.asked) != 0 {
		t.Errorf("asked the host with auto-upgrade off: %v", lab.asked)
	}
	if len(lab.threads) != 1 {
		t.Fatalf("threads: %v", lab.threads)
	}
	th := lab.threads[0]
	for _, want := range []string{"harbour: ", "1.48.0", "1.49.0", "flai host upgrade --version 1.49.0", "flai serve enable auto-upgrade"} {
		if !strings.Contains(th, want) {
			t.Errorf("the thread says %q:\n%s", want, th)
		}
	}
	for range 3 {
		lab.after(time.Minute)
		lab.look("harbour")
	}
	lab.restart()
	lab.after(time.Minute)
	lab.look("harbour")
	if len(lab.threads) != 1 {
		t.Fatalf("asked again about the same release, a restart included: %v", lab.threads)
	}
	lab.newest["harbour"] = "1.50.0"
	lab.after(time.Minute)
	lab.look("harbour")
	if len(lab.threads) != 2 || !strings.Contains(lab.threads[1], "1.50.0") {
		t.Errorf("a newer release asks again: %v", lab.threads)
	}
	if asks, err := lab.dir.upgradeAsks(); err != nil || asks["harbour"] != "1.50.0" {
		t.Errorf("the record holds the release asked about: %v %v", asks, err)
	}
}

// ADR-0158: a flai serve run outside flai host has no host to ask, so with
// auto-upgrade on it asks the operator on the thread, saying why.
func TestOutsideAHostTheOperatorIsAskedInstead(t *testing.T) {
	lab := newUpgradeLab(t)
	lab.underNo = true
	lab.restart()
	lab.auto["harbour"] = true
	lab.newest["harbour"] = "1.49.0"
	lab.look("harbour")
	if len(lab.threads) != 1 || !strings.Contains(lab.threads[0], "outside flai host") || strings.Contains(lab.threads[0], "flai serve enable auto-upgrade") {
		t.Errorf("the operator is asked, told why: %v", lab.threads)
	}
}

// ADR-0158: a development build is never behind, as FlaiOutdated holds, so
// it asks nothing and tells no one.
func TestADevelopmentBuildAsksNothing(t *testing.T) {
	lab := newUpgradeLab(t)
	lab.upgrader.outdated = func(root string) *release.Outdated { return release.FlaiOutdated(nil, root, "dev") }
	lab.auto["harbour"] = true
	lab.look("harbour")
	lab.after(time.Minute)
	lab.look("harbour")
	if len(lab.asked) != 0 || len(lab.threads) != 0 {
		t.Errorf("a dev build asked %v and told %v", lab.asked, lab.threads)
	}
}
