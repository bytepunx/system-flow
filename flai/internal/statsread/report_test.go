package statsread

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/metrics"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

var update = flag.Bool("update", false, "rewrite testdata/good-report.json from what Read and Compute give")

// noGit is a runner on a host without git. testdata/good sits inside this
// repository's checkout, so the real runner would read the monorepo's own
// history into Commits, which changes with every commit; without git,
// Commits stays nil and claims.drift is left out, whatever the history.
type noGit struct{ execx.System }

func (noGit) LookPath(name string) (string, error) {
	return "", errors.New(name + " is not on PATH in this test")
}

// reportCases are the windows and options the dashboard asks stats.get for,
// each read and computed over ../metrics/testdata/good at a fixed now.
var reportCases = []struct {
	name string
	opt  metrics.Options
}{
	{"story-1d", metrics.Options{Since: 24 * time.Hour}},
	{"story-7d", metrics.Options{Since: 7 * 24 * time.Hour}},
	{"story-30d", metrics.Options{Since: 30 * 24 * time.Hour}},
	{"story-365d", metrics.Options{Since: 365 * 24 * time.Hour}},
	{"story-30d-by-nature", metrics.Options{Since: 30 * 24 * time.Hour, By: "nature"}},
	{"story-7d-bucket-hour", metrics.Options{Since: 7 * 24 * time.Hour, Bucket: "hour"}},
}

// report reads the fixture with Read and computes the report for opt, as
// stats.get does, at the fixed now the metrics tests use.
func report(t *testing.T, opt metrics.Options) *metrics.Report {
	t.Helper()
	repo, err := workitem.Open("../metrics/testdata/good")
	if err != nil {
		t.Fatal(err)
	}
	items, read, err := Read(noGit{}, repo, slog.New(slog.DiscardHandler))
	if err != nil {
		t.Fatal(err)
	}
	read.Now = time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	read.Since, read.Type, read.By, read.Bucket = opt.Since, opt.Type, opt.By, opt.Bucket
	return metrics.Compute(items, read)
}

// S-0607: the report Read and Compute give over testdata/good stays what the
// serial code gave, whatever runs at once. The golden file was written by
// the serial code; run go test -update to rewrite it once a change to the
// report is intended, and review its diff. It holds one case a line, each
// report compact, since the 365-day one alone is about 250 KB indented.
func TestReportOverTheGoodFixtureIsTheGolden(t *testing.T) {
	var golden bytes.Buffer
	golden.WriteString("{\n")
	got := map[string][]byte{}
	for i, c := range reportCases {
		a, err := json.Marshal(report(t, c.opt))
		if err != nil {
			t.Fatal(err)
		}
		b, err := json.Marshal(report(t, c.opt))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(a, b) {
			t.Errorf("%s: two reads and computes of the same fixture differ at %s", c.name, firstDifference(a, b))
		}
		got[c.name] = a
		name, _ := json.Marshal(c.name)
		golden.Write(name)
		golden.WriteString(": ")
		golden.Write(a)
		if i < len(reportCases)-1 {
			golden.WriteString(",")
		}
		golden.WriteString("\n")
	}
	golden.WriteString("}\n")
	path := filepath.Join("testdata", "good-report.json")
	data, err := os.ReadFile(path)
	if *update || errors.Is(err, fs.ErrNotExist) {
		if err := os.WriteFile(path, golden.Bytes(), 0o644); err != nil {
			t.Fatal(err)
		}
		if !*update {
			t.Fatalf("%s was missing and is written now; review it and run the test again", path)
		}
		return
	}
	if err != nil {
		t.Fatal(err)
	}
	var want map[string]json.RawMessage
	if err := json.Unmarshal(data, &want); err != nil {
		t.Fatalf("%s: %v; run go test -update to write it again", path, err)
	}
	if len(want) != len(reportCases) {
		t.Errorf("%s holds %d cases, the test %d; run go test -update and review the diff", path, len(want), len(reportCases))
	}
	for _, c := range reportCases {
		if w, ok := want[c.name]; !ok {
			t.Errorf("%s holds no case %s", path, c.name)
		} else if !bytes.Equal(got[c.name], w) {
			t.Errorf("%s: the report over testdata/good differs from %s at %s; if the change is intended, run go test -update and review the diff", c.name, path, firstDifference(got[c.name], w))
		}
	}
}

// firstDifference shows where got first differs from want, with some of
// each around it.
func firstDifference(got, want []byte) string {
	i := 0
	for i < len(got) && i < len(want) && got[i] == want[i] {
		i++
	}
	around := func(b []byte) string {
		lo, hi := max(0, i-80), min(len(b), i+80)
		return string(b[lo:hi])
	}
	return "byte " + strconv.Itoa(i) + "\n got: …" + around(got) + "…\nwant: …" + around(want) + "…"
}
