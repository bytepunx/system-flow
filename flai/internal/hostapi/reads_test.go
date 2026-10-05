package hostapi

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
)

// readPhase is, per read a command also prints, a valid call and the phase
// its work is timed under (S-0159).
var readPhase = map[string]struct{ params, phase string }{
	"item.show":         {`{"id":"S-0001"}`, "item.show"},
	"item.move.preview": {`{"id":"E-0001"}`, "cancel.preview"},
	"accept.preview":    {`{"id":"S-0001"}`, "accept.preview"},
	"stream.diff":       {`{"id":"S-0001"}`, "stream.diff"},
	"stats.get":         {`{"since":"12w","type":"story","by":"nature","bucket":"week"}`, "stats.compute"},
	"publish.preview":   {`{}`, "release.pending"},
	// S-0217
	"order.by":           {`{"policy":"wsjf"}`, "order.by"},
	"promote.candidates": {`{"limit":3}`, "promote.candidates"},
	"release.evaluate":   {`{}`, "release.evaluate"},
}

// readRefused is, per read, params refused before anything is read.
var readRefused = map[string][]string{
	"item.show":         {`{"id":"--show"}`, `{"id":"../S-0001"}`},
	"item.move.preview": {`{"id":"../x"}`},
	"accept.preview":    {`{"id":"--no-push"}`},
	"stream.diff":       {`{"id":"../../etc"}`},
	"stats.get":         {`{"since":"30d; ls"}`, `{"type":"folder"}`, `{"by":"owner"}`, `{"bucket":"minute"}`, `{"bucket":"--json"}`},
	"publish.preview":   {`"--force"`},
	// S-0217
	"order.by":           {`{"policy":"random"}`, `{"policy":"--apply"}`, `{"policy":3}`},
	"promote.candidates": {`{"limit":-1}`, `{"limit":"3"}`},
	"release.evaluate":   {`"--apply"`, `[1]`},
}

func TestEveryReadIsCovered(t *testing.T) {
	reads := readMethods(time.Now)
	for name := range reads {
		if _, ok := readPhase[name]; !ok {
			t.Errorf("%s has no call in readPhase", name)
		}
		if len(readRefused[name]) == 0 {
			t.Errorf("%s has no refusals in readRefused", name)
		}
		if _, ok := specs()[name]; ok {
			t.Errorf("%s is still in the write table, which starts flai for it", name)
		}
	}
	if len(readPhase) != len(reads) || len(readRefused) != len(reads) {
		t.Errorf("readPhase names %d, readRefused %d, there are %d reads", len(readPhase), len(readRefused), len(reads))
	}
}

func TestAReadRefusesWhatIsNotDataBeforeReadingAnything(t *testing.T) {
	// a root that does not exist: anything past validation fails otherwise
	p := channel.Project{Key: "nowhere", Root: "/nonexistent/flai-test"}
	reads := readMethods(time.Now)
	for name, all := range readRefused {
		for _, params := range all {
			_, rerr := reads[name](context.Background(), p, json.RawMessage(params))
			if rerr == nil || rerr.Code != channel.CodeInvalidParams {
				t.Errorf("%s %s: %+v, want invalid params", name, params, rerr)
			}
		}
	}
}

func TestAReadIsAnsweredInThisProcess(t *testing.T) {
	p := harbour(t)
	table := Methods("test", func() time.Time { return t0.Add(time.Hour) })
	for name, call := range readPhase {
		got := phases(t, func(ctx context.Context) error {
			// the harbour is not a git repository: stream.diff and
			// publish.preview fail, as flai itself does there
			_, _ = table[name](ctx, p, json.RawMessage(call.params))
			return nil
		})
		if got[call.phase] != 1 {
			t.Errorf("%s: phase %s counted %d times, want once; phases %v", name, call.phase, got[call.phase], got)
		}
		for phase := range got {
			if strings.HasPrefix(phase, "exec.flai") {
				t.Errorf("%s started flai: phases %v", name, got)
			}
		}
	}
}

func TestAReadAnswersAsFlaiDid(t *testing.T) {
	p := harbour(t)
	table := Methods("test", func() time.Time { return t0.Add(time.Hour) })
	call := func(name, params string) (Written, *channel.Error) {
		res, rerr := table[name](context.Background(), p, json.RawMessage(params))
		w, _ := res.(Written)
		return w, rerr
	}

	w, rerr := call("item.show", `{"id":"S-1"}`)
	var shown struct{ ID, Hash string }
	if rerr != nil || json.Unmarshal(w.Data, &shown) != nil || shown.ID != "S-0001" || len(shown.Hash) != 64 || w.Warnings == nil {
		t.Errorf("item.show: %s %+v %+v", w.Data, w.Warnings, rerr)
	}

	// an ID that names nothing is the caller's mistake, as it was from flai
	if _, rerr := call("item.show", `{"id":"S-0099"}`); rerr == nil || rerr.Code != NotFound || !strings.Contains(rerr.Message, "S-0099 not found") {
		t.Errorf("item.show of nothing: %+v", rerr)
	}
	// a story not in review cannot be accepted: flai's failure, said as it said it
	if _, rerr := call("accept.preview", `{"id":"S-0001"}`); rerr == nil || rerr.Code != channel.CodeInternal || rerr.Message != "S-0001 is ready; a story is accepted from review" {
		t.Errorf("accept.preview of a ready story: %+v", rerr)
	}

	w, rerr = call("item.move.preview", `{"id":"E-0001"}`)
	var cancel struct {
		DryRun    bool `json:"dry_run"`
		Cancelled []struct{ ID string }
	}
	if rerr != nil || json.Unmarshal(w.Data, &cancel) != nil || !cancel.DryRun || len(cancel.Cancelled) == 0 {
		t.Errorf("item.move.preview: %s %+v", w.Data, rerr)
	}

	// S-0217: no policy is the manifest's, fifo when it names none, and the
	// order is computed, never applied
	w, rerr = call("order.by", `{}`)
	var order struct {
		Policy  string
		Applied bool
		Stories []struct{ ID string }
	}
	if rerr != nil || json.Unmarshal(w.Data, &order) != nil || order.Policy != "fifo" || order.Applied || order.Stories == nil {
		t.Errorf("order.by: %s %+v", w.Data, rerr)
	}
	w, rerr = call("promote.candidates", `{}`)
	var promote struct {
		Policy string
		Others []struct{ ID string }
	}
	if rerr != nil || json.Unmarshal(w.Data, &promote) != nil || promote.Policy != "fifo" {
		t.Errorf("promote.candidates: %s %+v", w.Data, rerr)
	}
}
