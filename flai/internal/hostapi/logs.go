package hostapi

// A log flai wrote, read for the dashboard (S-0785, ADR-0300): log.read, a
// read, answers at most 64 KB of a log from an offset, as images.log and
// kluster.log answer theirs, so that a path a status or an error names can
// be opened on the dashboard rather than in a shell on the host. The
// dashboard reaches the project only through flai (ADR-0029, ADR-0031), and
// flai's state folder holds its tokens beside the logs, so log.read reads a
// regular .log file only, by its real path, and only where flai writes the
// project's logs: under the project's .flai-cache but not in its worktrees,
// and in flai serve's state folder, Host.ServeDir, serve.log, the host's
// host.log beside it, and the agent runs, image jobs, and checks runs whose
// names or records say they are the project's.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The folders and files of flai serve's state folder that log.read reads, as
// flai serve names them: serve.Dir.Log, the launcher's openLog under agents,
// the image jobs' folder, and the checks runs' folder; and host.Dir.Log in
// the host's folder beside it, both named for the folder of the config file.
const (
	serveLogName   = "serve.log"
	serveAgents    = "agents"
	serveImageJobs = "image-jobs"
	serveChecks    = "checks"
	hostFolder     = "host"
	hostLogName    = "host.log"
)

// logSuffix ends the name of every file log.read reads.
const logSuffix = ".log"

// logReadable says what log.read reads, in each of its refusals.
const logReadable = "log.read reads only a .log file flai wrote for this project: one under the project's .flai-cache, but not in .flai-cache/worktrees, or in flai serve's state folder its serve.log, the host's host.log, and this project's agent runs, image jobs, and checks runs. Give the path a status or an error of this project names: absolute, ~/ for the home of the user flai runs as, or relative to the project's root"

var (
	// agentLogName is an agent run's log's name, <key>-<subject>-<stamp>.log.
	agentLogName = regexp.MustCompile(`^(.+)-\d{8}T\d{6}Z\.log$`)
	// checksLogName is a checks run's log's name, <story>-<stamp>.log.
	checksLogName = regexp.MustCompile(`^(S-\d{3,})-\d{8}T\d{6}Z\.log$`)
)

// LogRead is log.read's answer: a part of a log, read from Offset, with
// where to read from next, as images.log and kluster.log answer theirs.
type LogRead struct {
	// Path is the path asked, as it was given.
	Path string `json:"path"`
	// Exists says whether the log is there: one not written yet is none.
	Exists bool `json:"exists"`
	// Offset is where Text starts: the offset asked, or 0 after a reset.
	Offset int64 `json:"offset"`
	// Next is the offset to ask from next: Offset and Text's length.
	Next int64 `json:"next"`
	// Size is the log's size when it was read.
	Size int64 `json:"size"`
	// Text is at most 64 KB of the log from Offset, ending on a whole
	// character.
	Text string `json:"text"`
	// Reset is true when the log was shorter than the offset asked, and
	// Text then starts at 0.
	Reset bool `json:"reset,omitempty"`
}

// logReadParams are log.read's params.
type logReadParams struct {
	Path   string `json:"path"`
	Offset int64  `json:"offset"`
}

// logReadIn reads and checks log.read's params.
func logReadIn(raw json.RawMessage) (logReadParams, *channel.Error) {
	in, e := decode[logReadParams](raw)
	if e != nil {
		return in, e
	}
	switch {
	case strings.TrimSpace(in.Path) == "":
		return in, bad("give the path of the log to read, as a status or an error names it, such as .flai-cache/kluster/up.log")
	case strings.ContainsAny(in.Path, "\x00\n\r"):
		return in, bad("%q holds a line break or NUL, so it is no log's path; give the path as a status or an error names it", in.Path)
	case in.Offset < 0:
		return in, bad("%d is not an offset; give 0 to read from the start, or the next of the last read", in.Offset)
	}
	return in, nil
}

// logRefused is log.read's refusal of the path given, with why and what it
// reads.
func logRefused(given, why string) *channel.Error {
	return bad("%s cannot be read: %s. %s", given, why, logReadable)
}

// logReadDo answers log.read: at most 64 KB of a log flai wrote for the
// project, from an offset, whatever the host actions say.
func logReadDo(_ context.Context, host Host, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
	in, e := logReadIn(raw)
	if e != nil {
		return nil, e
	}
	path, e := logAbsolute(p.Root, in.Path)
	if e != nil {
		return nil, e
	}
	real, err := realPath(path)
	if err != nil {
		return nil, logRefused(in.Path, "its folder cannot be followed ("+klusterRelative(p.Root, err.Error())+")")
	}
	if e := logConfined(host, p, in.Path, real); e != nil {
		return nil, e
	}
	out, e := readConfinedLog(in.Path, real, in.Offset)
	if e != nil {
		return nil, e
	}
	out.Path = in.Path
	return out, nil
}

// logAbsolute is the path given as an absolute, clean path: ~/ is the home
// of the user flai runs as, and a relative path is from the project's root.
func logAbsolute(root, given string) (string, *channel.Error) {
	switch {
	case strings.HasPrefix(given, "~/"):
		home, err := os.UserHomeDir()
		if err != nil {
			return "", failed(fmt.Errorf("find the home folder to read %s from: %w; give the log's path absolute instead", given, err))
		}
		return filepath.Join(home, given[2:]), nil
	case filepath.IsAbs(given):
		return filepath.Clean(given), nil
	}
	return filepath.Join(root, given), nil
}

// realPath is path with every symbolic link in it followed: its own, or,
// where it or a folder of it is not there, its nearest folder that is,
// followed, with the rest of path after it.
func realPath(path string) (string, error) {
	rest := ""
	for dir := path; ; {
		real, err := filepath.EvalSymlinks(dir)
		if err == nil {
			return filepath.Join(real, rest), nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return "", err
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return path, nil
		}
		rest = filepath.Join(filepath.Base(dir), rest)
		dir = parent
	}
}

// below is real's path from base, slash separated, and whether real lies
// below base.
func below(base, real string) (string, bool) {
	rel, err := filepath.Rel(base, real)
	if err != nil || rel == "." || !filepath.IsLocal(rel) {
		return "", false
	}
	return filepath.ToSlash(rel), true
}

// logConfined refuses real, the real path of the log given, unless it is a
// .log file in one of the folders flai writes the project's logs to.
func logConfined(host Host, p channel.Project, given, real string) *channel.Error {
	if !strings.HasSuffix(filepath.Base(real), logSuffix) {
		return logRefused(given, "it is not a .log file, or it is a link to one that is not")
	}
	if host.ServeDir != "" {
		serve, err := realPath(host.ServeDir)
		if err != nil {
			return failed(fmt.Errorf("follow flai serve's state folder to read %s: %w; check that it can be read", given, err))
		}
		if rel, ok := below(serve, real); ok {
			return serveLogOwned(host, p, given, real, rel)
		}
		hostLog, err := realPath(filepath.Join(filepath.Dir(host.ServeDir), hostFolder, hostLogName))
		if err == nil && real == hostLog {
			return nil
		}
	}
	root, err := realPath(p.Root)
	if err != nil {
		return failed(fmt.Errorf("follow the project's root to read %s: %w; check that it can be read", given, err))
	}
	rel, ok := below(filepath.Join(root, ".flai-cache"), real)
	if !ok {
		return logRefused(given, "it lies outside the project's .flai-cache and flai serve's state folder, or a link in it leads there")
	}
	switch first, _, _ := strings.Cut(rel, "/"); first {
	case "worktrees":
		return logRefused(given, "it lies in a story's worktree, .flai-cache/worktrees, which is the story's branch, not a log flai keeps")
	case "serve", hostFolder:
		return logRefused(given, "it lies in the state folder of a flai whose config file is in the project's .flai-cache, which is not the flai serve answering")
	}
	return nil
}

// serveLogOwned refuses real, at rel in flai serve's state folder, unless
// it is serve.log or a log of the project's own: an agent run's, named for
// the project's key and a subject of the project's; an image job's that
// Host.Images says is the project's; or a checks run's of a story whose
// checks record names the project's root.
func serveLogOwned(host Host, p channel.Project, given, real, rel string) *channel.Error {
	folder, name, inFolder := strings.Cut(rel, "/")
	if !inFolder {
		if rel == serveLogName {
			return nil
		}
		return logRefused(given, "flai serve's state folder holds no other log log.read reads")
	}
	if strings.Contains(name, "/") {
		return logRefused(given, "flai serve keeps no log log.read reads that deep in its state folder")
	}
	switch folder {
	case serveAgents:
		if agentLogOwned(p, name) {
			return nil
		}
		return logRefused(given, "it is not an agent run of this project's: its name does not start with the project's key, "+p.Key+", and a story or a role of the project")
	case serveImageJobs:
		if host.Images == nil {
			return logRefused(given, "this flai runs no image jobs, so it cannot tell whose job it is")
		}
		_, path, ok := host.Images.JobLog(p.Root, strings.TrimSuffix(name, logSuffix))
		if job, err := realPath(path); ok && err == nil && job == real {
			return nil
		}
		return logRefused(given, "it is not the log of an image job of this project's; images.get lists them")
	case serveChecks:
		if checksLogOwned(p.Root, filepath.Dir(real), name) {
			return nil
		}
		return logRefused(given, "it is not a checks run of this project's: its story's checks record names another project, or none")
	}
	return logRefused(given, "flai serve keeps no log log.read reads in its "+folder+" folder")
}

// agentLogOwned says whether name, a log's name in flai serve's agents
// folder, is an agent run of the project p: <key>-<subject>-<stamp>.log,
// with p's key and, as its subject, a story's ID, a built-in strategic
// agent's kind, or a role the project defines. The subject's check keeps a
// project whose key is another's and a hyphen, such as shipyard-ui beside
// shipyard, from reading the other's runs.
func agentLogOwned(p channel.Project, name string) bool {
	m := agentLogName.FindStringSubmatch(name)
	if m == nil || p.Key == "" {
		return false
	}
	subject, ok := strings.CutPrefix(m[1], p.Key+"-")
	switch {
	case !ok:
		return false
	case storyID.MatchString(subject), slices.Contains(workitem.ActivityKinds, subject):
		return true
	}
	return slices.Contains(roleNames(p.Root), subject)
}

// roleNames are the names of the roles the project at root defines; none
// when its manifest or its role files cannot be read.
func roleNames(root string) []string {
	m, err := manifest.Load(filepath.Join(root, manifest.File))
	if err != nil {
		return nil
	}
	all, _, err := roles.Load(root, m.Layout["design"])
	if err != nil {
		return nil
	}
	out := make([]string, 0, len(all))
	for _, r := range all {
		out = append(out, r.Name)
	}
	return out
}

// checksLogOwned says whether name, a log's name in flai serve's checks
// folder dir, is a checks run of the project at root: <story>-<stamp>.log,
// whose story's checks record, <story>.json beside it, names root. The
// record is the story's newest run's, so a story's older runs are read
// while its newest is the project's.
func checksLogOwned(root, dir, name string) bool {
	m := checksLogName.FindStringSubmatch(name)
	if m == nil {
		return false
	}
	data, err := os.ReadFile(filepath.Join(dir, m[1]+".json"))
	if err != nil {
		return false
	}
	var rec struct {
		Root string `json:"root"`
	}
	if json.Unmarshal(data, &rec) != nil || rec.Root == "" {
		return false
	}
	ran, err1 := realPath(rec.Root)
	own, err2 := realPath(root)
	return err1 == nil && err2 == nil && ran == own
}

// readConfinedLog reads at most 64 KB of the log at real from offset, a log
// not there being none, read from 0. It refuses what is not a regular file,
// and a file that real no longer names, or no longer names alone, once it is
// open, so that a link put in its place between the check and the read is
// not followed.
func readConfinedLog(given, real string, offset int64) (LogRead, *channel.Error) {
	before, err := os.Lstat(real)
	if errors.Is(err, fs.ErrNotExist) {
		return absentLog(offset), nil
	}
	if err != nil {
		return LogRead{}, failed(fmt.Errorf("look at the log %s: %w; check that flai can read it", given, err))
	}
	if !before.Mode().IsRegular() {
		return LogRead{}, logRefused(given, "it is not a regular file but a link, a folder, or a device")
	}
	f, err := os.Open(real)
	if err != nil {
		return LogRead{}, failed(fmt.Errorf("open the log %s: %w; check that flai can read it", given, err))
	}
	defer func() { _ = f.Close() }()
	st, err := f.Stat()
	if err != nil {
		return LogRead{}, failed(fmt.Errorf("look at the log %s: %w; check that flai can read it", given, err))
	}
	after, err := os.Lstat(real)
	again, err2 := filepath.EvalSymlinks(real)
	if err != nil || err2 != nil || again != real || !os.SameFile(before, st) || !os.SameFile(after, st) {
		return LogRead{}, logRefused(given, "it changed while it was opened, to another file or a link")
	}
	out, err := readOpenLog(f, st.Size(), offset)
	if err != nil {
		return LogRead{}, failed(fmt.Errorf("read the log %s: %w; check that flai can read it", given, err))
	}
	return out, nil
}

// absentLog is the part of a log not there read from offset: none, from 0,
// reset when offset was past 0.
func absentLog(offset int64) LogRead {
	return LogRead{Reset: offset > 0}
}

// readOpenLog reads at most imageLogChunk bytes of f, of size bytes, from
// offset, ending on a whole character, or from 0 when offset is past its end.
func readOpenLog(f io.ReaderAt, size, offset int64) (LogRead, error) {
	out := LogRead{Exists: true, Size: size, Offset: offset}
	if offset > size {
		out.Offset, out.Reset = 0, true
	}
	buf := make([]byte, min(size-out.Offset, imageLogChunk))
	n, err := f.ReadAt(buf, out.Offset)
	if err != nil && !errors.Is(err, io.EOF) {
		return LogRead{}, err
	}
	buf = wholeChars(buf[:n])
	out.Text, out.Next = string(buf), out.Offset+int64(len(buf))
	return out, nil
}

// logSpecs are the log methods: log.read, a read answered whatever the host
// actions say, whose build checks the params and whose do reads the log.
func logSpecs() map[string]spec {
	return map[string]spec{
		// log.read: at most 64 KB of a log flai wrote for the project, from
		// an offset, the dashboard's log viewer.
		"log.read": {reads: true, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			_, e := logReadIn(raw)
			return nil, "", e
		}, do: logReadDo},
	}
}
