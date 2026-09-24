package serve

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
)

// S-0109: what one flai serve records of a detached write, the next reads.
func TestRequestsOutliveTheServeThatWroteThem(t *testing.T) {
	d := Dir(t.TempDir())
	at := time.Date(2026, 9, 24, 7, 0, 0, 0, time.UTC)
	res := &hostapi.Written{Data: json.RawMessage(`{"restarting":true}`), Warnings: []string{}}
	if err := d.Requests(func(all map[string]hostapi.Request) {
		all["k"] = hostapi.Request{Started: at, PID: 42, Until: at.Add(time.Hour), Done: true, Result: res}
	}); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(d.requests())
	if err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("requests.json: %v %v", info, err)
	}
	var got hostapi.Request
	if err := d.Requests(func(all map[string]hostapi.Request) { got = all["k"] }); err != nil {
		t.Fatal(err)
	}
	var data bytes.Buffer
	if got.Result != nil {
		_ = json.Compact(&data, got.Result.Data)
	}
	if !got.Started.Equal(at) || got.PID != 42 || !got.Done || data.String() != `{"restarting":true}` {
		t.Errorf("read back %+v", got)
	}
}

// A file that is not the records is an error that names it, not a silent
// fresh start that would do repeats again.
func TestUnreadableRequestsSayWhichFile(t *testing.T) {
	d := Dir(t.TempDir())
	for _, content := range []string{"{", "null"} {
		if err := os.WriteFile(d.requests(), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
		err := d.Requests(func(all map[string]hostapi.Request) { all["k"] = hostapi.Request{} })
		switch content {
		case "{":
			if err == nil || !strings.Contains(err.Error(), "requests.json") {
				t.Errorf("%q: %v", content, err)
			}
		default:
			if err != nil {
				t.Errorf("%q: %v", content, err)
			}
		}
	}
}
