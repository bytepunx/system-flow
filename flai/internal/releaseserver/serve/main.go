// Command serve answers GitHub's release endpoints from GoReleaser dist
// folders on 127.0.0.1, for the smoke tier's install and self-upgrade checks
// (S-0340). It prints its base URL as one line on stdout once it listens,
// logs to stderr, and stops on SIGINT or SIGTERM, or once the process
// -owner-pid names is gone. scripts/release-server.sh runs it.
//
//	serve -repo local/flai -dir flai/dist [-dir other/dist] [-addr 127.0.0.1:0] [-owner-pid N]
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"

	"github.com/bytepunx/system-flow/flai/internal/logx"
	"github.com/bytepunx/system-flow/flai/internal/releaseserver"
)

// dirs is the repeatable -dir flag.
type dirs []string

func (d *dirs) String() string { return strings.Join(*d, ",") }

func (d *dirs) Set(v string) error {
	*d = append(*d, v)
	return nil
}

func main() {
	log := logx.New(os.Stderr, func() logx.Options {
		o := logx.FromEnv()
		o.IsTerminal = term.IsTerminal(int(os.Stderr.Fd()))
		return o
	}())
	if err := run(log); err != nil {
		logx.Fatal(log, "release server stopped", "component", "releaseserver", "err", err.Error())
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	var folders dirs
	repo := flag.String("repo", "local/flai", "repository name the releases are served under, owner/name")
	addr := flag.String("addr", "127.0.0.1:0", "address to listen on; port 0 picks a free one")
	owner := flag.Int("owner-pid", 0, "stop once the process with this PID is gone; 0 runs until SIGINT or SIGTERM")
	flag.Var(&folders, "dir", "GoReleaser dist folder to serve releases from (repeatable)")
	flag.Parse()
	if flag.NArg() > 0 {
		return fmt.Errorf("unexpected arguments %q: give each dist folder with -dir", flag.Args())
	}
	if strings.Count(*repo, "/") != 1 || strings.HasPrefix(*repo, "/") || strings.HasSuffix(*repo, "/") {
		return fmt.Errorf("-repo %q is not owner/name, such as local/flai", *repo)
	}
	releases, err := releaseserver.Scan(folders...)
	if err != nil {
		return err
	}
	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w; give -addr a free address, or 127.0.0.1:0", *addr, err)
	}
	srv := &http.Server{Handler: logRequests(log, releaseserver.New(*repo, releases)), ReadHeaderTimeout: 10 * time.Second}
	tags := make([]string, len(releases))
	for i, r := range releases {
		tags[i] = r.Tag
	}
	base := "http://" + ln.Addr().String()
	log.Info("release server listening", "component", "releaseserver", "url", base, "repo", *repo, "releases", strings.Join(tags, ","), "owner_pid", *owner)
	if _, err := fmt.Fprintln(os.Stdout, base); err != nil {
		_ = ln.Close()
		return fmt.Errorf("print the base URL: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	served := make(chan error, 1)
	go func() { served <- srv.Serve(ln) }()
	ownerGone := watchOwner(ctx, *owner)
	select {
	case err := <-served:
		return fmt.Errorf("serve on %s: %w", base, err)
	case <-ctx.Done():
		log.Info("release server stopping", "component", "releaseserver", "reason", "signal")
	case <-ownerGone:
		log.Info("release server stopping", "component", "releaseserver", "reason", "owner gone", "owner_pid", *owner)
	}
	shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdown); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("shut down the server on %s: %w", base, err)
	}
	return nil
}

// watchOwner closes the channel it answers once the process pid is gone,
// checking each second; with pid 0 it never closes.
func watchOwner(ctx context.Context, pid int) <-chan struct{} {
	gone := make(chan struct{})
	if pid <= 0 {
		return gone
	}
	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()
		for alive(pid) {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
			}
		}
		close(gone)
	}()
	return gone
}

// alive reports whether the process pid still runs. Unix answers signal 0;
// Windows finds no process that has exited, and sends no signal 0.
func alive(pid int) bool {
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	if runtime.GOOS == "windows" {
		_ = p.Release()
		return true
	}
	err = p.Signal(syscall.Signal(0))
	return err == nil || errors.Is(err, syscall.EPERM)
}

// status records the status a handler writes.
type status struct {
	http.ResponseWriter
	code int
}

func (s *status) WriteHeader(code int) {
	s.code = code
	s.ResponseWriter.WriteHeader(code)
}

// logRequests logs one line per request: method, path, status, duration.
func logRequests(log *slog.Logger, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &status{ResponseWriter: w, code: http.StatusOK}
		h.ServeHTTP(sw, r)
		log.Info("request served", "component", "releaseserver", "method", r.Method, "path", r.URL.EscapedPath(), "status", sw.code, "duration_ms", time.Since(start).Milliseconds())
	})
}
