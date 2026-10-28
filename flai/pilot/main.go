// Command pilot is the entrypoint of flai's agent container images
// (ADR-0212): it reads one start request on standard input, starts the inner
// harness it names through its adapter, relays the harness's output, forwards
// SIGTERM and SIGINT to it, and exits with its exit code.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bytepunx/system-flow/flai/internal/pilot"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	os.Exit(pilot.Main(context.Background(), os.Stdin, os.Stdout, os.Stderr, os.Environ(), signals))
}
