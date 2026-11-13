package cmd

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestProfilesAreOfferedOnLoopbackOnly(t *testing.T) {
	var logs bytes.Buffer
	log := slog.New(slog.NewTextHandler(&logs, nil))
	t.Setenv(pprofEnv, "0.0.0.0:0")
	servePprof(context.Background(), log)
	if !strings.Contains(logs.String(), "not a loopback address") {
		t.Errorf("an address beyond this machine: %s", logs.String())
	}

	// the kernel picks the port as servePprof listens, and the log names it: a
	// port the test closed and handed over could be taken first (I-0372)
	logs.Reset()
	t.Setenv(pprofEnv, "127.0.0.1:0")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	servePprof(ctx, log)
	_, url, found := strings.Cut(logs.String(), "url=")
	url, _, _ = strings.Cut(url, " ")
	url = strings.Trim(strings.TrimSpace(url), `"`)
	if !found || !strings.HasPrefix(url, "http://127.0.0.1:") || strings.HasPrefix(url, "http://127.0.0.1:0/") {
		t.Fatalf("the log names the address profiles are offered at: %s", logs.String())
	}
	var res *http.Response
	var err error
	for range 50 {
		if res, err = http.Get(url + "cmdline"); err == nil {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if err != nil {
		t.Fatal(err)
	}
	_ = res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("status %d", res.StatusCode)
	}
}
