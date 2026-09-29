package cmd

import (
	"bytes"
	"context"
	"log/slog"
	"net"
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

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	_ = ln.Close()
	t.Setenv(pprofEnv, addr)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	servePprof(ctx, log)
	var res *http.Response
	for range 50 {
		if res, err = http.Get("http://" + addr + "/debug/pprof/cmdline"); err == nil {
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
