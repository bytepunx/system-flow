package hostapi

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// S-0788, ADR-0302: accept.resolve runs flai accept --resolve as the
// project's owner with every host action off, since it merges nothing and
// reaches no forge, and answers its send back as written; its refusal of a
// story nothing holds, exit 4, is Refused with the refused object as data.
func TestAcceptResolveRunsUngatedAndItsRefusalIsData(t *testing.T) {
	p := withDocs(t) // owner: olive
	off := Host{Enabled: func(string, string) bool { return false }}
	params := json.RawMessage(`{"id":"S-0001",` + rid + `}`)

	sent := `{"id":"S-0001","reason":"S-0001 was sent back to in-progress","sent_back":true,"branch":"story/S-0001","base":"main","conflicts":[],"thread":"TH-0001","blockers":["b"]}`
	rec := &recorder{ran: Ran{Stdout: []byte(sent)}}
	res, e := writeMethods(rec.run, time.Now, off)["accept.resolve"](context.Background(), p, params)
	if e != nil {
		t.Fatalf("accept.resolve must run with every host action off: %+v", e)
	}
	if len(rec.runs) != 1 || strings.Join(rec.runs[0].Args, " ") != "accept S-0001 --resolve --by=olive --json" {
		t.Fatalf("accept.resolve ran %+v", rec.runs)
	}
	if w, ok := res.(Written); !ok || string(w.Data) != sent {
		t.Errorf("the send back is the answer: %#v", res)
	}

	reason := "nothing holds S-0001's acceptance, so its agent has nothing to resolve: accept it with flai accept S-0001"
	rec = &recorder{ran: Ran{Exit: 4, Stdout: []byte(`{"refused":{"id":"S-0001","reason":"` + reason + `"}}`),
		Events: []map[string]any{{"level": "FATAL", "msg": "command failed", "err": reason}}}}
	_, e = writeMethods(rec.run, time.Now, off)["accept.resolve"](context.Background(), p, params)
	if e == nil || e.Code != Refused || e.Message != reason {
		t.Fatalf("exit 4 should be Refused with flai's reason: %+v", e)
	}
	if data, _ := e.Data.(map[string]any); data["id"] != "S-0001" || data["reason"] != reason {
		t.Errorf("the error's data: %+v", e.Data)
	}
}
