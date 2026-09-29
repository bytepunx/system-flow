package cmd

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"time"
)

// pprofEnv names the loopback address flai serve offers Go's profiles on
// (S-0152). Unset, it offers none.
const pprofEnv = "FLAI_PPROF_ADDR"

// servePprof offers CPU, heap, goroutine, and the other runtime profiles
// on the loopback address pprofEnv names, until ctx ends: go tool pprof
// http://127.0.0.1:6060/debug/pprof/profile?seconds=30. A profile names
// the host's paths, so an address beyond this machine is refused.
func servePprof(ctx context.Context, log *slog.Logger) {
	addr := os.Getenv(pprofEnv)
	if addr == "" {
		return
	}
	if !loopback(addr) {
		log.Warn("profiles not offered: the address is not a loopback address", "component", "serve", "env", pprofEnv, "addr", addr)
		return
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Warn("profiles not offered: cannot listen", "component", "serve", "env", pprofEnv, "addr", addr, "err", err.Error())
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	srv := &http.Server{Handler: mux, ReadHeaderTimeout: 10 * time.Second}
	go func() {
		<-ctx.Done()
		_ = srv.Close()
	}()
	log.Info("profiles offered", "component", "serve", "url", "http://"+ln.Addr().String()+"/debug/pprof/")
	go func() {
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Warn("profiles no longer offered", "component", "serve", "err", err.Error())
		}
	}()
}
