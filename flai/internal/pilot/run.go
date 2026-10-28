package pilot

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/bytepunx/system-flow/flai/internal/harness"
)

// sessionName is a session pilot can record in its cache as one file name.
var sessionName = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

// Main is pilot's run (S-0507, ADR-0212): it reads the start request on
// stdin, checks it against environ, pilot's own environment, and runs it,
// forwarding each signal that arrives on signals to the inner harness and
// SIGTERM when ctx is done. It returns pilot's exit code: the inner harness's,
// or ExitUsage for a request it refuses, after one line on stderr.
func Main(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer, environ []string, signals <-chan os.Signal) int {
	req, err := Read(stdin)
	if err == nil {
		err = Check(req, getenv(environ), isDir)
	}
	if err != nil {
		say(stderr, err.Error())
		return ExitUsage
	}
	return Run(ctx, req, stdout, stderr, environ, signals)
}

// Run starts the inner harness req, a request Check passed, names through
// its adapter, in a new session or resuming the one the cache records, and
// relays its standard
// output to stdout unchanged; the harness writes its standard error to
// stderr itself. The child's environment is environ, then the adapter's
// variables, then those the adapter copies from environ, then FLAI_AGENT,
// FLAI_SESSION, and FLAI_STORY from the request, a later one of a name
// winning, as os/exec keeps the last. It returns the harness's exit code, or
// 128 and the signal's number when a signal ended it, and ExitUsage, after
// one line on stderr, when the harness cannot be started.
func Run(ctx context.Context, req harness.Request, stdout, stderr io.Writer, environ []string, signals <-chan os.Signal) int {
	cmd, marker, err := prepare(&req, stderr, environ)
	var out io.Reader
	if err == nil {
		out, err = start(cmd, req)
	}
	if err != nil {
		say(stderr, err.Error())
		return ExitUsage
	}
	return finish(ctx, cmd, out, req, marker, stdout, stderr, signals)
}

// prepare makes the command that runs req's inner harness, and the marker
// that records its session: it sets req.Resume when the cache records the
// session already, writes the files the adapter's start needs, and builds
// the child's environment. It starts nothing.
func prepare(req *harness.Request, stderr io.Writer, environ []string) (*exec.Cmd, string, error) {
	marker, err := resume(req)
	if err != nil {
		return nil, "", err
	}
	adapter, ok := harness.Adapters[req.Inner]
	if !ok {
		return nil, "", fmt.Errorf("the start request's inner %s is not a harness pilot can run; set inner to one of %s", strconv.Quote(req.Inner), strings.Join(innerNames(), ", "))
	}
	spec, err := adapter.Start(*req, harness.Host{})
	if err != nil {
		return nil, "", fmt.Errorf("the inner harness %s cannot start the request: %w", req.Inner, err)
	}
	if len(spec.Argv) == 0 {
		return nil, "", fmt.Errorf("the inner harness %s gave no program to run; its adapter in harness.Adapters must name one", req.Inner)
	}
	for _, f := range spec.Files {
		if err := write(f); err != nil {
			return nil, "", err
		}
	}
	env, err := childEnv(*req, spec, environ)
	if err != nil {
		return nil, "", err
	}
	path, err := exec.LookPath(spec.Argv[0])
	if err != nil {
		return nil, "", fmt.Errorf("the inner harness %s's program %s is not found on pilot's PATH: %w; build the agent image over the harness's layer, which installs it", req.Inner, strconv.Quote(spec.Argv[0]), err)
	}
	cmd := &exec.Cmd{Path: path, Args: spec.Argv, Env: env, Stderr: stderr}
	return cmd, marker, nil
}

// resume is the file under the cache that records req's session, made ready
// to write: none when req names no session. It sets req.Resume when the file
// is there already, so that the harness goes on in the session. It refuses a
// session that is not one file's name.
func resume(req *harness.Request) (string, error) {
	if req.Session == "" {
		return "", nil
	}
	if !sessionName.MatchString(req.Session) || req.Session == "." || req.Session == ".." {
		return "", errors.New("the start request's Session is not a session's name; name it with letters, digits, '.', '_', and '-' alone")
	}
	dir := filepath.Join(req.Cache, "pilot", "sessions")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", fmt.Errorf("making the folder %s that records the sessions pilot started: %w; mount the cache writable", dir, err)
	}
	marker := filepath.Join(dir, req.Session)
	if _, err := os.Stat(marker); err == nil {
		req.Resume = true
	}
	return marker, nil
}

// write writes f with its mode, making its folder, and sets the mode of a
// file that was there already.
func write(f harness.File) error {
	if err := os.MkdirAll(filepath.Dir(f.Path), 0o700); err != nil {
		return fmt.Errorf("making the folder of %s, which the inner harness reads: %w; mount the cache writable", f.Path, err)
	}
	if err := os.WriteFile(f.Path, f.Data, f.Mode); err != nil {
		return fmt.Errorf("writing %s, which the inner harness reads: %w; mount the cache writable", f.Path, err)
	}
	if err := os.Chmod(f.Path, f.Mode); err != nil {
		return fmt.Errorf("setting the mode of %s, which the inner harness reads: %w; mount the cache writable", f.Path, err)
	}
	return nil
}

// childEnv is the inner harness's environment, as Run says. It refuses a
// variable the adapter copies from one that is unset or empty in environ,
// naming the variable, never a value.
func childEnv(req harness.Request, spec harness.Start, environ []string) ([]string, error) {
	get := getenv(environ)
	env := append(append([]string{}, environ...), spec.Env...)
	for child, parent := range spec.EnvFrom {
		value := get(parent)
		if value == "" {
			return nil, fmt.Errorf("the inner harness %s reads %s from the variable %s, which is unset or empty in pilot's environment; pass it to the container", req.Inner, child, strconv.Quote(parent))
		}
		env = append(env, child+"="+value)
	}
	for _, v := range []struct{ name, value string }{{"FLAI_AGENT", req.Name}, {"FLAI_SESSION", req.Session}, {"FLAI_STORY", req.Story}} {
		if v.value != "" {
			env = append(env, v.name+"="+v.value)
		}
	}
	return env, nil
}

// start starts cmd and gives the pipe its standard output is on, saying so
// on its standard error, without the session's name.
func start(cmd *exec.Cmd, req harness.Request) (io.Reader, error) {
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("making the pipe the inner harness %s writes to: %w", req.Inner, err)
	}
	how := "in a new session"
	if req.Resume {
		how = "resuming its session"
	}
	say(cmd.Stderr, "starting "+req.Inner+" "+how)
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("starting the inner harness %s: %w", req.Inner, err)
	}
	return out, nil
}

// finish relays out, cmd's standard output, to stdout and forwards signals
// to cmd until it ends, records its session in marker, and returns its exit
// code.
func finish(ctx context.Context, cmd *exec.Cmd, out io.Reader, req harness.Request, marker string, stdout, stderr io.Writer, signals <-chan os.Signal) int {
	if marker != "" {
		if err := os.WriteFile(marker, nil, 0o600); err != nil {
			say(stderr, fmt.Sprintf("recording the session in %s failed, so a later start will not resume it: %v", marker, err))
		}
	}
	done := make(chan struct{})
	forwarded := make(chan struct{})
	go func() {
		forward(ctx, cmd.Process, signals, done)
		close(forwarded)
	}()
	rerr := relay(stdout, out)
	werr := cmd.Wait()
	close(done)
	<-forwarded
	if rerr != nil {
		say(stderr, fmt.Sprintf("relaying the inner harness %s's output failed: %v", req.Inner, rerr))
	}
	var exit *exec.ExitError
	if werr != nil && !errors.As(werr, &exit) {
		say(stderr, fmt.Sprintf("waiting for the inner harness %s failed: %v", req.Inner, werr))
	}
	code := exitCode(cmd.ProcessState)
	say(stderr, fmt.Sprintf("%s exited with code %d", req.Inner, code))
	return code
}

// forward sends p each signal that arrives on signals, and SIGTERM once ctx
// is done, until done is closed.
func forward(ctx context.Context, p *os.Process, signals <-chan os.Signal, done <-chan struct{}) {
	stop := ctx.Done()
	for {
		select {
		case <-done:
			return
		case sig := <-signals:
			// The child may have ended already; Wait reports how.
			_ = p.Signal(sig)
		case <-stop:
			_ = p.Signal(syscall.SIGTERM)
			stop = nil
		}
	}
}

// relay copies r to w unchanged, writing each line as it arrives, and a
// line longer than the reader's buffer in parts as they arrive, so that its
// memory stays bounded. After w fails it reads r to its end, so that the
// child never blocks on a full pipe, and returns the failure.
func relay(w io.Writer, r io.Reader) error {
	br := bufio.NewReaderSize(r, 64<<10)
	var werr error
	for {
		chunk, err := br.ReadSlice('\n')
		if len(chunk) > 0 && werr == nil {
			_, werr = w.Write(chunk)
		}
		switch {
		case errors.Is(err, bufio.ErrBufferFull):
			continue
		case errors.Is(err, io.EOF):
			return werr
		case err != nil:
			return errors.Join(werr, err)
		}
	}
}

// exitCode is the exit code of the child that state describes: its own, or
// 128 and the signal's number when a signal ended it; ExitUsage when it
// never ran to an end.
func exitCode(state *os.ProcessState) int {
	if state == nil {
		return ExitUsage
	}
	if ws, ok := state.Sys().(syscall.WaitStatus); ok && ws.Signaled() {
		return 128 + int(ws.Signal())
	}
	return state.ExitCode()
}

// getenv looks a variable up in environ, the last of its name winning; empty
// when it is unset.
func getenv(environ []string) func(string) string {
	return func(name string) string {
		value := ""
		for _, kv := range environ {
			if k, v, ok := strings.Cut(kv, "="); ok && k == name {
				value = v
			}
		}
		return value
	}
}

// isDir is whether path is a folder.
func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// say writes one line of pilot's own on w, prefixed pilot:.
func say(w io.Writer, line string) {
	_, _ = fmt.Fprintln(w, "pilot: "+strings.ReplaceAll(line, "\n", " "))
}
