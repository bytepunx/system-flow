package mcpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/taskdone"
)

// installedFlai is the flai binary at the path an MCP server was started
// from, which closes a task in place of the server's own flai once that is
// outdated (ADR-0160): a stdio server keeps its process until its session
// ends, while a host upgrade replaces the file at its path.
type installedFlai struct {
	path string
	dir  string   // where it runs: the project's main checkout
	env  []string // the server's environment, its agent and session included
}

// version is the installed flai's version, as flai version --json prints it.
func (f installedFlai) version(ctx context.Context) (string, error) {
	stdout, stderr, err := f.run(ctx, "version", "--json")
	if err != nil {
		return "", fmt.Errorf("%s version --json: %w%s", f.path, err, said("said", stderr))
	}
	var v struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(stdout, &v); err != nil || strings.TrimSpace(v.Version) == "" {
		return "", fmt.Errorf("%s version --json printed no version%s", f.path, said("printed", string(stdout)))
	}
	return strings.TrimSpace(v.Version), nil
}

// taskDone closes the task with flai task done --json, as the server's agent
// and session, and answers the result it prints. A close that stops at a
// step exits non-zero and still prints its result, which is the answer;
// output that is not a result is the error, with what the command printed.
func (f installedFlai) taskDone(ctx context.Context, in TaskDoneIn) (taskdone.Result, error) {
	args := []string{"task", "done", in.Task, "--json"}
	if in.Message != "" {
		args = append(args, "-m", in.Message)
	}
	if in.Log != "" {
		args = append(args, "--log", in.Log)
	}
	stdout, stderr, err := f.run(ctx, args...)
	var res taskdone.Result
	if jerr := json.Unmarshal(stdout, &res); jerr == nil && res.Task != "" {
		return res, nil
	}
	failed := ""
	if err != nil {
		failed = ": " + err.Error()
	}
	return taskdone.Result{Task: in.Task}, fmt.Errorf("the installed flai printed no result closing %s: %s task done %s --json%s%s%s",
		in.Task, f.path, in.Task, failed, said("printed", string(stdout)), said("said", stderr))
}

// run runs the installed flai with args in its folder and environment, and
// answers its standard output and standard error apart: the result is on
// the first, log lines on the second.
func (f installedFlai) run(ctx context.Context, args ...string) (stdout []byte, stderr string, err error) {
	cmd := exec.CommandContext(ctx, f.path, args...)
	cmd.Dir, cmd.Env = f.dir, f.env
	var out, errOut bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errOut
	err = cmd.Run()
	var exit *exec.ExitError
	if err != nil && !errors.As(err, &exit) {
		err = fmt.Errorf("cannot run it: %w", err)
	}
	return out.Bytes(), errOut.String(), err
}

// said is "; it <verb> <the end of out, quoted>", for an error message, or
// nothing when out is blank.
func said(verb, out string) string {
	if strings.TrimSpace(out) == "" {
		return ""
	}
	return fmt.Sprintf("; it %s %q", verb, tail(out))
}

// tail is s trimmed and cut to its last 400 bytes, where a failure is said.
func tail(s string) string {
	const most = 400
	s = strings.TrimSpace(s)
	if len(s) > most {
		return "..." + s[len(s)-most:]
	}
	return s
}
