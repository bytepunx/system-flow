package metrics_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/metrics"
	"github.com/bytepunx/system-flow/flai/internal/statsread"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// statsGolden is where flai stats --json over testdata/good is pinned.
var statsGolden = filepath.Join("testdata", "good.stats.golden.json")

// statsNow is the time every pinned run reads as now.
var statsNow = time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)

// statsRun is one flai stats --json run the golden pins, by its flags.
type statsRun struct {
	name, typ, since, by, bucket string
}

// statsRuns are the pinned runs, in the golden's order.
var statsRuns = []statsRun{
	{name: "story 30d", typ: workitem.Story, since: metrics.DefaultWindow, bucket: metrics.DefaultBucket},
	{name: "story 30d by nature", typ: workitem.Story, since: metrics.DefaultWindow, by: "nature", bucket: metrics.DefaultBucket},
	{name: "story 30d by parent", typ: workitem.Story, since: metrics.DefaultWindow, by: "parent", bucket: metrics.BucketWeek},
	{name: "task 90d", typ: workitem.Task, since: "90d", bucket: metrics.DefaultBucket},
}

// noGit is a runner on a host without git, so that the stories' commits are
// not read from the repository testdata/good sits in, and claims leave the
// drift out.
type noGit struct{}

var errNoGit = errors.New("git is not installed")

func (noGit) Run(string, string, ...string) (string, error)              { return "", errNoGit }
func (noGit) RunInput(string, string, string, ...string) (string, error) { return "", errNoGit }
func (noGit) LookPath(string) (string, error)                            { return "", errNoGit }

// statsJSON is what flai stats --json prints for each of statsRuns over
// testdata/good, as cmd's stats builds it: statsread.Read's items and options,
// the run's flags, and the report encoded with two spaces of indent and a
// newline. No serve journal is read, as on a host without one.
func statsJSON(t *testing.T) [][]byte {
	t.Helper()
	out := make([][]byte, len(statsRuns))
	for i, run := range statsRuns {
		repo, err := workitem.Open(filepath.Join("testdata", "good"))
		if err != nil {
			t.Fatal(err)
		}
		window, err := metrics.ParseWindow(run.since)
		if err != nil {
			t.Fatal(err)
		}
		items, opt, err := statsread.Read(noGit{}, repo, slog.New(slog.DiscardHandler))
		if err != nil {
			t.Fatal(err)
		}
		opt.Now, opt.Since, opt.Type, opt.By, opt.Bucket = statsNow, window, run.typ, run.by, run.bucket
		var b bytes.Buffer
		enc := json.NewEncoder(&b)
		enc.SetIndent("", "  ")
		if err := enc.Encode(metrics.Compute(items, opt)); err != nil {
			t.Fatal(err)
		}
		out[i] = b.Bytes()
	}
	return out
}

// S-0750, ADR-0305: on the baseline workflow flai stats --json gives over
// testdata/good the same bytes as before its figures became readings of the
// workflow's timers. The golden is the runs' output one after the other, and
// was written by main's code before S-0750; S-0750 must not write it again,
// since a difference is what this test is there to find.
func TestStatsJSONOnGoodIsTheGolden(t *testing.T) {
	want, err := os.ReadFile(statsGolden)
	if err != nil {
		t.Fatal(err)
	}
	got := statsJSON(t)
	dec := json.NewDecoder(bytes.NewReader(want))
	for i, run := range statsRuns {
		var pinned json.RawMessage
		if err := dec.Decode(&pinned); err != nil {
			t.Fatalf("%s holds no report for run %s: %v", statsGolden, run.name, err)
		}
		if g := bytes.TrimSuffix(got[i], []byte("\n")); !bytes.Equal(g, pinned) {
			t.Errorf("%s: flai stats --json differs from %s at %s", run.name, statsGolden, differAt(g, pinned))
		}
	}
	if all := bytes.Join(got, nil); !bytes.Equal(all, want) {
		t.Errorf("the runs' output differs from %s at %s", statsGolden, differAt(all, want))
	}
}

// differAt shows where got first differs from want, with some of each
// around it.
func differAt(got, want []byte) string {
	i := 0
	for i < len(got) && i < len(want) && got[i] == want[i] {
		i++
	}
	around := func(b []byte) string { return string(b[max(0, i-120):min(len(b), i+120)]) }
	return "byte " + strconv.Itoa(i) + "\n got: …" + around(got) + "…\nwant: …" + around(want) + "…"
}
