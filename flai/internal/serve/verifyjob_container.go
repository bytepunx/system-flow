package serve

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/closeout"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The close-out of a story whose agent runs on the container harness
// (ADR-0274, I-0295). Its worktree is the container's to write, so the host
// runs nothing from it: the verify, which runs the worktree's scripts and the
// test tiers they select, runs in a container from the story's image with
// the story's mounts (harness.ContainerCommand), and the steps the host must
// take, the commit and the sync, are flai's own code and git. The host's git
// takes its git dir from the main checkout's record of the worktree, never
// from the files the container can write, and runs the main checkout's hooks
// alone, with no fsmonitor. It mirrors scripts/close-out.sh's steps and its
// last line, at most closeOutRounds rounds (ADR-0156).

// closeOutRounds is how many times a close-out syncs and verifies again when
// main moved while it ran (ADR-0156).
const closeOutRounds = 3

// containerAgent is a story's agent that runs on the container harness and
// the host's say about agents, which the close-out's image is decided from.
type containerAgent struct {
	agent *manifest.Agent
	cfg   AgentConfig
}

// containerAgentOf is story's agent when it runs on the container harness,
// resolved as flai serve starts it (ADR-0251, ADR-0261); nil when it runs on
// any other or the story cannot be read.
func (p *verifyPass) containerAgentOf(repo *workitem.Repo, story string) *containerAgent {
	cfg := AgentConfig{}
	if p.o.Agent != nil {
		cfg = p.o.Agent(p.e.Root)
	}
	it, err := repo.Get(story)
	if err != nil {
		return nil
	}
	a := cfg.storyAgent(repo, it)
	if a == nil || a.Harness != harness.Container {
		return nil
	}
	return &containerAgent{agent: a, cfg: cfg}
}

// containerCloseOut is one close-out job of a container story.
type containerCloseOut struct {
	p        *verifyPass
	agent    containerAgent
	story    string
	worktree string
	req      closeout.Request
	id       string // the job's start, which names its containers
}

// String is the close-out's command as the job's start logs it.
func (c *containerCloseOut) String() string {
	return "flai verify " + c.story + " --record-issues in the story's container, then the commit and the sync on the host"
}

// verifyArgv is the image's flai verify of the story, with all and fresh as
// the request asks.
func (c *containerCloseOut) verifyArgv() []string {
	argv := []string{harness.ImageFlai, "verify", c.story, "--record-issues"}
	if c.req.All {
		argv = append(argv, "--all")
	}
	if c.req.Fresh {
		argv = append(argv, "--fresh")
	}
	return argv
}

// run is the job's body: the verify in the container, the commit, the sync
// check, and, when main moved, flai stream sync and the next round, each
// written to out, which ends with one line naming the story, the outcome,
// the step it stopped at, and the round, as the script's last line does. It
// answers the exit code of the step it stopped at, 0 when every step passed,
// and an error when a step could not run; nothing from the worktree runs on
// the host, and a story whose container cannot run fails rather than run on
// the host (ADR-0274 §5).
func (c *containerCloseOut) run(ctx context.Context, out io.Writer) (int, error) {
	p, root, story := c.p, c.p.e.Root, c.story
	round, step := 1, "the container check"
	end := func(code int, err error) (int, error) {
		in := ""
		if round > 1 {
			in = fmt.Sprintf(" in round %d of %d", round, closeOutRounds)
		}
		switch {
		case err != nil:
			fmt.Fprintf(out, "close-out: %s stopped at %s%s (%v)\n", story, step, in, err)
			return -1, err
		case code == 0:
			fmt.Fprintf(out, "close-out: %s passed every step%s; ready to move to review\n", story, in)
		default:
			fmt.Fprintf(out, "close-out: %s stopped at %s%s (exit %d)\n", story, step, in, code)
		}
		return code, nil
	}
	host, image, err := c.image(ctx, out)
	if err == nil && image == "" {
		err = fmt.Errorf("the container harness has no image for the inner harness %s of %s's agent: name one with flai serve agent harness container --image %s=<image>", innerOf(c.agent.agent), story, innerOf(c.agent.agent))
	}
	if err != nil {
		return end(-1, fmt.Errorf("%s's agent runs in a container, so its close-out runs there and nothing runs on the host: %w", story, err))
	}
	flai, err := p.exe()
	if err != nil {
		return end(-1, fmt.Errorf("this flai's executable is not known, so the host's sync check cannot run: %w", err))
	}
	base := closeOutEnv(withoutGitEnv(os.Environ()), c.req.By)
	hostEnv := append(slices.Clone(base), hostGitConfig(root)...)
	gitEnv := append(slices.Clone(hostEnv), storyGitDir(root, story, c.worktree)...)
	for {
		step = "flai verify"
		name := c.id
		if round > 1 {
			name = fmt.Sprintf("%s-r%d", c.id, round)
		}
		cmd, err := harness.ContainerCommand(harness.CommandRun{Root: root, Story: story, Purpose: "closeout", ID: name, Image: image,
			Env: []string{"FLAI_AGENT"}, Argv: c.verifyArgv()}, host)
		if err != nil {
			return end(-1, err)
		}
		fmt.Fprintf(out, "close-out: flai verify %s --record-issues in the container %s\n", story, cmd.Name)
		if code, err := p.run(ctx, c.worktree, append(slices.Clone(base), cmd.Env...), out, cmd.Argv); err != nil || code != 0 {
			return end(code, err)
		}

		step = "commit"
		if err := gitPointers(root, story, c.worktree); err != nil {
			return end(-1, err)
		}
		if code, err := c.commit(ctx, out, gitEnv); err != nil || code != 0 {
			return end(code, err)
		}

		step = "the sync check"
		fmt.Fprintf(out, "close-out: flai verify %s --sync-only\n", story)
		var said bytes.Buffer
		code, err := p.run(ctx, c.worktree, hostEnv, io.MultiWriter(out, &said), []string{flai, "verify", story, "--sync-only"})
		if err != nil {
			return end(code, err)
		}
		if code == 0 {
			break
		}
		// Main moved, such as by an acceptance, when it stopped at the sync
		// step; at any other step, or after the last round, it stops here.
		if !strings.HasPrefix(lastLine(said.String()), "verify: "+story+" stopped at sync (") || round >= closeOutRounds {
			return end(code, nil)
		}
		round++
		step = "flai stream sync"
		fmt.Fprintf(out, "close-out: main moved during the run; flai stream sync %s (round %d of %d)\n", story, round, closeOutRounds)
		if err := gitPointers(root, story, c.worktree); err != nil {
			return end(-1, err)
		}
		if code, err := p.run(ctx, c.worktree, hostEnv, out, []string{flai, "stream", "sync", story}); err != nil || code != 0 {
			return end(code, err)
		}
	}

	step = "the clean worktree check"
	var status bytes.Buffer
	if code, err := p.run(ctx, c.worktree, gitEnv, &status, []string{"git", "status", "--porcelain"}); err != nil || code != 0 {
		return end(code, err)
	}
	if strings.TrimSpace(status.String()) != "" {
		fmt.Fprintf(out, "%sclose-out: the worktree is not clean after the commit\n", status.String())
		return end(1, nil)
	}
	return end(0, nil)
}

// image is the container harness's host and the image the story's close-out
// runs in. When its agent images are missing or stale and flai serve runs
// image jobs, it asks for the job that makes them, under the project's main
// checkout, so that it joins a start's, writes to out that it waits for it,
// and waits, then judges the image again (ADR-0295 §5): the close-out's
// verdict, and its verified event, come after the job. A failed job is its
// error, which names its cause, its step, and its log. With no image jobs a
// missing image refuses, naming flai image ensure.
func (c *containerCloseOut) image(ctx context.Context, out io.Writer) (harness.Host, string, error) {
	p, root, docker := c.p, c.p.e.Root, dockerOf(c.p.o)
	host, image, wait, err := containerImageAsked(docker, p.o.imageJobs, p.e.Key, root, c.agent.cfg, c.agent.agent)
	if err != nil || wait == nil {
		return host, image, err
	}
	j := wait.job
	fmt.Fprintf(out, "close-out: %s: waiting for flai serve's image job %s, which makes them, its log %s\n", wait.why, j.ID(), j.Log())
	p.o.Logger.Info("close-out waits for an image job", "component", "serve", "root", root, "story", c.story, "job_id", c.story+"-"+c.id, "image_job", j.ID(), "image", wait.image)
	if _, err := j.wait(ctx); err != nil {
		return harness.Host{}, "", err
	}
	fmt.Fprintf(out, "close-out: the image job %s passed\n", j.ID())
	return containerImage(docker, root, c.agent.cfg, c.agent.agent)
}

// commit commits what is outstanding in the worktree on the host with git
// and the request's message, as the script's commit step does; with
// something outstanding and no message it says so and answers 1.
func (c *containerCloseOut) commit(ctx context.Context, out io.Writer, gitEnv []string) (int, error) {
	run := c.p.run
	var status bytes.Buffer
	if code, err := run(ctx, c.worktree, gitEnv, &status, []string{"git", "status", "--porcelain"}); err != nil || code != 0 {
		fmt.Fprint(out, status.String())
		return code, err
	}
	if strings.TrimSpace(status.String()) == "" {
		return 0, nil
	}
	if c.req.Message == "" {
		fmt.Fprintf(out, "%sclose-out: uncommitted changes and no message; ask for the close-out again with the commit message\n", status.String())
		return 1, nil
	}
	fmt.Fprintln(out, "close-out: commit")
	if code, err := run(ctx, c.worktree, gitEnv, out, []string{"git", "add", "-A"}); err != nil || code != 0 {
		return code, err
	}
	return run(ctx, c.worktree, gitEnv, out, []string{"git", "commit", "-m", c.req.Message})
}

// hostGitConfig is the configuration every git the host runs for a container
// story's close-out is given through the environment, which outranks every
// configuration file: the main checkout's hooks and no fsmonitor, so a
// core.hooksPath or core.fsmonitor the container wrote into a file git reads
// runs nothing on the host (ADR-0274 §3).
func hostGitConfig(root string) []string {
	return []string{"GIT_CONFIG_COUNT=2",
		"GIT_CONFIG_KEY_0=core.hooksPath", "GIT_CONFIG_VALUE_0=" + filepath.Join(root, ".git", "hooks"),
		"GIT_CONFIG_KEY_1=core.fsmonitor", "GIT_CONFIG_VALUE_1=false"}
}

// storyGitDir names story's git dir, its common dir, and its work tree from
// the main checkout at root, never from the worktree's .git file or the
// registration's commondir, which the container can write. It is given to
// the git commands the close-out runs itself; the flai it runs on the host
// reaches git in more than one checkout, so it gets hostGitConfig alone, and
// gitPointers checks the files before each of its steps.
func storyGitDir(root, story, worktree string) []string {
	return []string{"GIT_DIR=" + filepath.Join(root, ".git", "worktrees", story), "GIT_COMMON_DIR=" + filepath.Join(root, ".git"), "GIT_WORK_TREE=" + worktree}
}

// gitPointers refuses when the worktree's .git file no longer names story's
// registration in the main checkout at root, or the registration's commondir
// no longer names the main checkout's .git: a git on the host would then read
// a repository, its configuration and hooks among them, that the container
// chose.
func gitPointers(root, story, worktree string) error {
	reg := filepath.Join(root, ".git", "worktrees", story)
	common := filepath.Join(root, ".git")
	refuse := func(what string) error {
		return fmt.Errorf("%s, and only %s's container could have changed it, so the host runs no git there: run git worktree repair %s in the main checkout %s, look at what the container wrote, and ask again", what, story, worktree, root)
	}
	dotGit := filepath.Join(worktree, ".git")
	data, err := os.ReadFile(dotGit)
	if err != nil {
		return refuse(fmt.Sprintf("%s cannot be read (%v)", dotGit, err))
	}
	line, _, _ := strings.Cut(string(data), "\n")
	gitdir, ok := strings.CutPrefix(strings.TrimSpace(line), "gitdir:")
	if !ok || !samePath(resolvedFrom(worktree, strings.TrimSpace(gitdir)), reg) {
		return refuse(fmt.Sprintf("%s names %q, not the registration %s", dotGit, strings.TrimSpace(line), reg))
	}
	file := filepath.Join(reg, "commondir")
	data, err = os.ReadFile(file)
	if err != nil {
		return refuse(fmt.Sprintf("%s cannot be read (%v)", file, err))
	}
	if named := strings.TrimSpace(string(data)); !samePath(resolvedFrom(reg, named), common) {
		return refuse(fmt.Sprintf("%s names %q, not the main checkout's %s", file, named, common))
	}
	return nil
}

// resolvedFrom is path, made absolute against dir when it is relative.
func resolvedFrom(dir, path string) string {
	if path == "" || filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(dir, path)
}

// samePath says whether a and b name the same folder, through links.
func samePath(a, b string) bool {
	if filepath.Clean(a) == filepath.Clean(b) {
		return true
	}
	ra, errA := filepath.EvalSymlinks(a)
	rb, errB := filepath.EvalSymlinks(b)
	return errA == nil && errB == nil && ra == rb
}

// withoutGitEnv is env without the variables that point git at a repository
// or give it configuration, so the close-out's own say is the only one.
func withoutGitEnv(env []string) []string {
	return slices.DeleteFunc(slices.Clone(env), func(kv string) bool {
		name, _, _ := strings.Cut(kv, "=")
		switch name {
		case "GIT_DIR", "GIT_COMMON_DIR", "GIT_WORK_TREE", "GIT_INDEX_FILE", "GIT_CONFIG_COUNT", "GIT_CONFIG_PARAMETERS":
			return true
		}
		return strings.HasPrefix(name, "GIT_CONFIG_KEY_") || strings.HasPrefix(name, "GIT_CONFIG_VALUE_")
	})
}

// lastLine is the last line of s that is not empty.
func lastLine(s string) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	return strings.TrimSpace(lines[len(lines)-1])
}
