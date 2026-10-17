package components

// Serving a story's worktree (S-0442): a worktree is a fresh checkout, so
// its sub-project's dependencies may not be installed. Install installs
// them by the lock file the branch carries, as the main checkout's were, and
// ServeWorktree installs them before it serves.

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// InstallTimeout is how long Install lets a package manager install a
// sub-project's dependencies when InstallOptions gives no Timeout.
const InstallTimeout = 10 * time.Minute

// ErrNoLockFile is wrapped by the error Install returns when a sub-project
// has neither its dependencies installed nor a lock file to install them by.
var ErrNoLockFile = errors.New("no lock file")

// lockFile is a lock file a package manager installs by, and the command
// that installs exactly what it locks.
type lockFile struct {
	name    string
	command string
	args    []string
}

// lockFiles are the lock files Install knows, in the order it prefers them
// when a folder holds more than one.
var lockFiles = []lockFile{
	{name: "pnpm-lock.yaml", command: "pnpm", args: []string{"install", "--frozen-lockfile"}},
	{name: "package-lock.json", command: "npm", args: []string{"ci"}},
	{name: "yarn.lock", command: "yarn", args: []string{"install", "--frozen-lockfile"}},
}

// commandLine is the install command as a shell would read it.
func (l lockFile) commandLine() string {
	return strings.Join(append([]string{l.command}, l.args...), " ")
}

// InstallOptions says whose dependencies Install installs, and how.
type InstallOptions struct {
	// Root is the checkout's root, absolute: the highest folder a lock file
	// is looked for in, for a workspace's lock file.
	Root string
	// Dir is the sub-project's folder, absolute, at or below Root.
	Dir string
	// Runner runs the package manager; System when nil.
	Runner Runner
	// Logger gets the install's start and end; nothing is logged when nil.
	Logger *slog.Logger
	// Timeout is how long the package manager may take; zero is
	// InstallTimeout.
	Timeout time.Duration
}

// Install installs the dependencies of the sub-project in opts.Dir when its
// vite is not installed, by the first lock file found in its folder or a
// folder above it up to opts.Root: pnpm install --frozen-lockfile for
// pnpm-lock.yaml, npm ci for package-lock.json, yarn install
// --frozen-lockfile for yarn.lock, run in the lock file's folder. It reports
// whether it ran a package manager. A failed install's error carries the
// command and its last lines of output; no lock file is ErrNoLockFile, and
// vite still missing after the install is ErrViteMissing.
func Install(ctx context.Context, opts InstallOptions) (bool, error) {
	if !filepath.IsAbs(opts.Root) || !filepath.IsAbs(opts.Dir) {
		return false, fmt.Errorf("cannot install the dependencies of the sub-project in %q: flai needs it and the checkout's root (got %q) as absolute paths", opts.Dir, opts.Root)
	}
	opts.Root, opts.Dir = filepath.Clean(opts.Root), filepath.Clean(opts.Dir)
	if rel, err := filepath.Rel(opts.Root, opts.Dir); err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return false, fmt.Errorf("cannot install the dependencies of the sub-project in %s: it is not inside the checkout at %s; check its path in system-flow.yaml", opts.Dir, opts.Root)
	}
	if info, err := os.Stat(opts.Dir); errors.Is(err, fs.ErrNotExist) || err == nil && !info.IsDir() {
		return false, fmt.Errorf("cannot install the dependencies of the sub-project in %s: its folder is not there; check its path in system-flow.yaml, and that the branch has it", opts.Dir)
	} else if err != nil {
		return false, fmt.Errorf("cannot install the dependencies of the sub-project in %s: cannot read its folder: %w", opts.Dir, err)
	}
	if _, ok, err := findVite(opts.Dir); err != nil || ok {
		return false, err
	}
	lock, at, err := findLockFile(opts.Root, opts.Dir)
	if err != nil {
		return false, err
	}
	if opts.Runner == nil {
		opts.Runner = System{}
	}
	if opts.Logger == nil {
		opts.Logger = slog.New(slog.DiscardHandler)
	}
	if opts.Timeout <= 0 {
		opts.Timeout = InstallTimeout
	}
	if _, err := opts.Runner.LookPath(lock.command); err != nil {
		return false, fmt.Errorf("cannot install the dependencies of the sub-project in %s: %s is not on PATH, and its lock file %s needs it; install %s, then run this again", opts.Dir, lock.command, filepath.Join(at, lock.name), lock.command)
	}
	log := opts.Logger.With("component", "components", "dir", at, "lock_file", lock.name, "command", lock.commandLine())
	log.Info("installing sub-project dependencies")
	start := time.Now()
	run, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()
	stdout, stderr, err := opts.Runner.Run(run, at, lock.command, "", lock.args...)
	if err != nil {
		out := &tail{}
		_, _ = out.Write(stdout)
		_, _ = out.Write([]byte("\n"))
		_, _ = out.Write(stderr)
		_, _ = out.Write([]byte("\n"))
		switch {
		case ctx.Err() != nil:
			return true, fmt.Errorf("stopped %s in %s before it installed the sub-project's dependencies: %w%s", lock.commandLine(), at, ctx.Err(), detail(out.last()))
		case errors.Is(run.Err(), context.DeadlineExceeded):
			return true, fmt.Errorf("%s in %s did not install the sub-project's dependencies within %s, so flai stopped it; run it there yourself, then run this again%s", lock.commandLine(), at, opts.Timeout, detail(out.last()))
		default:
			return true, fmt.Errorf("%s in %s failed: %w; fix what its last lines say, then run this again%s", lock.commandLine(), at, err, detail(out.last()))
		}
	}
	log.Info("sub-project dependencies installed", "duration", time.Since(start).Round(time.Millisecond).String())
	if _, ok, err := findVite(opts.Dir); err != nil {
		return true, err
	} else if !ok {
		return true, fmt.Errorf("%w in %s after %s in %s, and flai runs the sub-project's own vite from its node_modules/.bin: add vite to the sub-project's own devDependencies, then run this again", ErrViteMissing, opts.Dir, lock.commandLine(), at)
	}
	return true, nil
}

// findLockFile is the first lock file of lockFiles in dir, then in each
// folder above it up to root, and the folder that holds it. dir is at or
// below root.
func findLockFile(root, dir string) (lockFile, string, error) {
	for at := dir; ; at = filepath.Dir(at) {
		for _, l := range lockFiles {
			if _, err := os.Stat(filepath.Join(at, l.name)); err == nil {
				return l, at, nil
			} else if !errors.Is(err, fs.ErrNotExist) {
				return lockFile{}, "", fmt.Errorf("cannot look for a lock file at %s: %w; check that the checkout can be read", filepath.Join(at, l.name), err)
			}
		}
		if at == root || filepath.Dir(at) == at {
			break
		}
	}
	names := make([]string, len(lockFiles))
	for i, l := range lockFiles {
		names[i] = l.name
	}
	return lockFile{}, "", fmt.Errorf("cannot install the dependencies of the sub-project in %s: %w (%s) in it or a folder above it up to %s, and flai installs only what a lock file pins; commit the sub-project's lock file to the branch, or install its dependencies there yourself, then run this again", dir, ErrNoLockFile, strings.Join(names, ", "), root)
}

// WorktreeOptions says whose components ServeWorktree serves from a
// story's worktree, and how it installs their dependencies.
type WorktreeOptions struct {
	// ServeOptions are Serve's, with Root the worktree's root.
	ServeOptions
	// Runner runs the package manager; System when nil.
	Runner Runner
	// Logger gets the install's start and end; nothing is logged when nil.
	Logger *slog.Logger
	// InstallTimeout is how long the package manager may take; zero is
	// InstallTimeout.
	InstallTimeout time.Duration
}

// ServeWorktree is Serve in a story's worktree: it first installs the
// sub-project's dependencies, as Install does, when they are missing.
// Options Serve refuses are refused by Serve, without an install.
func ServeWorktree(ctx context.Context, opts WorktreeOptions) (*DevServer, error) {
	if _, ok := opts.Sub.ComponentSettings(); ok && filepath.IsAbs(opts.Root) && opts.Project != "" {
		_, err := Install(ctx, InstallOptions{
			Root:    opts.Root,
			Dir:     filepath.Join(opts.Root, filepath.FromSlash(opts.Sub.Path)),
			Runner:  opts.Runner,
			Logger:  opts.Logger,
			Timeout: opts.InstallTimeout,
		})
		if err != nil {
			return nil, fmt.Errorf("cannot serve sub-project %s's components: %w", opts.Sub.Name, err)
		}
	}
	return Serve(ctx, opts.ServeOptions)
}
