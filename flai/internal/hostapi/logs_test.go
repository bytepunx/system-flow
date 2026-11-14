package hostapi

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
)

// logsHome is a home folder with flai's state folder in it, .flai, holding
// its config file, and the host whose flai serve keeps its state there, with
// jobs as its image jobs.
func logsHome(t *testing.T, jobs ImageJobs) (string, Host) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	writeLog(t, filepath.Join(home, ".flai", "config.json"), `{"token":"secret"}`)
	return home, Host{ServeDir: filepath.Join(home, ".flai", "serve"), Images: jobs}
}

// writeLog writes text to path, making its folders.
func writeLog(t *testing.T, path, text string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}

// readLogPart asks log.read of p, on host, for path from offset.
func readLogPart(t *testing.T, host Host, p channel.Project, path string, offset int64) (LogRead, *channel.Error) {
	t.Helper()
	m, ok := MethodsFor("test", func() time.Time { return t0 }, host)["log.read"]
	if !ok {
		t.Fatal("no method log.read")
	}
	params, err := json.Marshal(map[string]any{"path": path, "offset": offset})
	if err != nil {
		t.Fatal(err)
	}
	res, e := m(context.Background(), p, params)
	if e != nil {
		return LogRead{}, e
	}
	data, err := json.Marshal(res)
	if err != nil {
		t.Fatal(err)
	}
	var out LogRead
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("log.read: %v\n%s", err, data)
	}
	return out, nil
}

func TestLogReadIsAReadUnderNoHostAction(t *testing.T) {
	if sp := logSpecs()["log.read"]; sp.action != "" || !sp.reads || sp.do == nil {
		t.Errorf("log.read is a read, answered whatever the host actions say: %+v", sp)
	}
}

func TestLogReadReadsEachKindOfLogFlaiWritesForTheProject(t *testing.T) {
	p := shipyard(t)
	jobs := &fakeImageJobs{state: ImagesState{Ended: []ImageJob{{ID: "20261113T093000Z-001", State: ImageJobEnded}}}, logs: map[string]string{}}
	home, host := logsHome(t, jobs)
	serve := host.ServeDir
	jobs.logs["20261113T093000Z-001"] = filepath.Join(serve, "image-jobs", "20261113T093000Z-001.log")
	writeLog(t, filepath.Join(serve, "checks", "S-0001.json"), `{"story":"S-0001","root":"`+p.Root+`"}`)
	writeLog(t, filepath.Join(p.Root, "design", "workflows", "roles", "reviewer.md"), reviewerRole)
	for name, c := range map[string]struct{ file, given string }{
		"kluster step, relative":   {filepath.Join(p.Root, ".flai-cache", "kluster", "logs", "story-s-0001", "deploy.log"), ".flai-cache/kluster/logs/story-s-0001/deploy.log"},
		"kluster up":               {filepath.Join(p.Root, ".flai-cache", "kluster", "up.log"), ".flai-cache/kluster/up.log"},
		"close-out, absolute":      {filepath.Join(p.Root, ".flai-cache", "verify", "jobs", "S-0001.log"), filepath.Join(p.Root, ".flai-cache", "verify", "jobs", "S-0001.log")},
		"merge queue tier":         {filepath.Join(p.Root, ".flai-cache", "merge-queue", "batches", "20261102-055211", "unit.log"), ".flai-cache/merge-queue/batches/20261102-055211/unit.log"},
		"story's agent run, ~/":    {filepath.Join(serve, "agents", "shipyard-S-0001-20261114T100000Z.log"), "~/.flai/serve/agents/shipyard-S-0001-20261114T100000Z.log"},
		"planner's run":            {filepath.Join(serve, "agents", "shipyard-planner-20261114T100000Z.log"), filepath.Join(serve, "agents", "shipyard-planner-20261114T100000Z.log")},
		"defined role's run":       {filepath.Join(serve, "agents", "shipyard-reviewer-20261114T100000Z.log"), filepath.Join(serve, "agents", "shipyard-reviewer-20261114T100000Z.log")},
		"image job":                {jobs.logs["20261113T093000Z-001"], "~/.flai/serve/image-jobs/20261113T093000Z-001.log"},
		"checks run":               {filepath.Join(serve, "checks", "S-0001-20261114T100000Z.log"), filepath.Join(serve, "checks", "S-0001-20261114T100000Z.log")},
		"flai serve's own log":     {filepath.Join(serve, "serve.log"), "~/.flai/serve/serve.log"},
		"the host's log beside it": {filepath.Join(home, ".flai", "host", "host.log"), "~/.flai/host/host.log"},
	} {
		writeLog(t, c.file, name+"\n")
		got, e := readLogPart(t, host, p, c.given, 0)
		if e != nil {
			t.Errorf("%s: %v", name, e)
			continue
		}
		want := LogRead{Path: c.given, Exists: true, Offset: 0, Next: int64(len(name) + 1), Size: int64(len(name) + 1), Text: name + "\n"}
		if got != want {
			t.Errorf("%s:\n got %+v\nwant %+v", name, got, want)
		}
	}
}

func TestLogReadReadsFromTheOffsetAndResetsPastTheEnd(t *testing.T) {
	p := shipyard(t)
	writeLog(t, filepath.Join(p.Root, ".flai-cache", "kluster", "up.log"), "hello world")
	got, e := readLogPart(t, Host{}, p, ".flai-cache/kluster/up.log", 6)
	if e != nil {
		t.Fatal(e)
	}
	if want := (LogRead{Path: ".flai-cache/kluster/up.log", Exists: true, Offset: 6, Next: 11, Size: 11, Text: "world"}); got != want {
		t.Errorf("from 6:\n got %+v\nwant %+v", got, want)
	}
	got, e = readLogPart(t, Host{}, p, ".flai-cache/kluster/up.log", 11)
	if e != nil || got.Text != "" || got.Next != 11 || got.Reset {
		t.Errorf("from the end: %+v %v, want no text, next 11, no reset", got, e)
	}
	got, e = readLogPart(t, Host{}, p, ".flai-cache/kluster/up.log", 40)
	if e != nil {
		t.Fatal(e)
	}
	if want := (LogRead{Path: ".flai-cache/kluster/up.log", Exists: true, Offset: 0, Next: 11, Size: 11, Text: "hello world", Reset: true}); got != want {
		t.Errorf("past the end:\n got %+v\nwant %+v", got, want)
	}
}

func TestLogReadAnswersALogNotWrittenYetAsNone(t *testing.T) {
	p := shipyard(t)
	got, e := readLogPart(t, Host{}, p, ".flai-cache/kluster/logs/story-s-0002/build.log", 5)
	if e != nil {
		t.Fatal(e)
	}
	if want := (LogRead{Path: ".flai-cache/kluster/logs/story-s-0002/build.log", Reset: true}); got != want {
		t.Errorf("a log not there:\n got %+v\nwant %+v", got, want)
	}
}

func TestLogReadAnswersAt64KBEndingOnAWholeCharacter(t *testing.T) {
	p := shipyard(t)
	big := strings.Repeat("a", 70000)
	writeLog(t, filepath.Join(p.Root, ".flai-cache", "big.log"), big)
	got, e := readLogPart(t, Host{}, p, ".flai-cache/big.log", 0)
	if e != nil {
		t.Fatal(e)
	}
	if len(got.Text) != 64<<10 || got.Next != 64<<10 || got.Size != 70000 {
		t.Errorf("a 70000-byte log: %d bytes, next %d, size %d, want 65536, 65536, 70000", len(got.Text), got.Next, got.Size)
	}
	got, e = readLogPart(t, Host{}, p, ".flai-cache/big.log", got.Next)
	if e != nil || len(got.Text) != 70000-64<<10 || got.Next != 70000 {
		t.Errorf("the rest: %d bytes, next %d, %v, want %d and 70000", len(got.Text), got.Next, e, 70000-64<<10)
	}
	writeLog(t, filepath.Join(p.Root, ".flai-cache", "cut.log"), strings.Repeat("a", 64<<10-1)+"é and on")
	got, e = readLogPart(t, Host{}, p, ".flai-cache/cut.log", 0)
	if e != nil {
		t.Fatal(e)
	}
	if len(got.Text) != 64<<10-1 || got.Next != 64<<10-1 || strings.ContainsRune(got.Text, '\uFFFD') {
		t.Errorf("a character across 64 KB: %d bytes, next %d, want both 65535 so that é is read next", len(got.Text), got.Next)
	}
	got, e = readLogPart(t, Host{}, p, ".flai-cache/cut.log", got.Next)
	if e != nil || got.Text != "é and on" {
		t.Errorf("the next read: %q %v, want it to start with é", got.Text, e)
	}
}

func TestLogReadRefusesWhatIsNotTheProjectsLog(t *testing.T) {
	p := shipyard(t)
	jobs := &fakeImageJobs{state: ImagesState{Ended: []ImageJob{{ID: "20261113T093000Z-001", State: ImageJobEnded}}},
		logs: map[string]string{"20261113T093000Z-001": "/elsewhere/20261113T093000Z-001.log"}}
	home, host := logsHome(t, jobs)
	serve := host.ServeDir
	outside := t.TempDir()
	writeLog(t, filepath.Join(outside, "secret.log"), "the content of a file outside")
	writeLog(t, filepath.Join(filepath.Dir(p.Root), "outside.log"), "the content of a file outside")
	writeLog(t, filepath.Join(p.Root, ".flai-cache", "config.json"), "{}")
	writeLog(t, filepath.Join(p.Root, ".flai-cache", "worktrees", "S-0001", ".flai-cache", "x.log"), "branch")
	writeLog(t, filepath.Join(p.Root, ".flai-cache", "serve", "agents", "other-S-0001-20261114T100000Z.log"), "a dev flai's")
	writeLog(t, filepath.Join(p.Root, "build.log"), "the project's own file")
	writeLog(t, filepath.Join(serve, "agents", "other-S-0001-20261114T100000Z.log"), "another project's")
	writeLog(t, filepath.Join(serve, "agents", "shipyard-ui-S-0001-20261114T100000Z.log"), "shipyard-ui's")
	writeLog(t, filepath.Join(serve, "image-jobs", "20261113T093000Z-001.log"), "another project's job")
	writeLog(t, filepath.Join(serve, "image-jobs", "20261113T093000Z-002.log"), "no job of the project's")
	writeLog(t, filepath.Join(serve, "checks", "S-0002.json"), `{"story":"S-0002","root":"/elsewhere"}`)
	writeLog(t, filepath.Join(serve, "checks", "S-0002-20261114T100000Z.log"), "another project's checks")
	writeLog(t, filepath.Join(serve, "journal.log"), "not a log flai serve writes")
	writeLog(t, filepath.Join(home, ".flai", "other.log"), "beside the config")
	links := map[string]string{
		filepath.Join(p.Root, ".flai-cache", "out.log"):                                 filepath.Join(outside, "secret.log"),
		filepath.Join(p.Root, ".flai-cache", "config.log"):                              filepath.Join(home, ".flai", "config.json"),
		filepath.Join(serve, "agents", "shipyard-S-0009-20261114T100000Z.log"):          filepath.Join(home, ".flai", "config.json"),
		filepath.Join(p.Root, ".flai-cache", "kluster"):                                 outside,
		filepath.Join(p.Root, ".flai-cache", "dangling.log"):                            filepath.Join(outside, "gone.log"),
		filepath.Join(serve, "agents", "shipyard-S-0008-20261114T100000Z.log"):          filepath.Join(serve, "agents", "other-S-0001-20261114T100000Z.log"),
		filepath.Join(p.Root, ".flai-cache", "merge-queue"):                             filepath.Join(p.Root, ".flai-cache", "worktrees"),
		filepath.Join(p.Root, ".flai-cache", "verify", "jobs", "S-0001-alias.log"):      filepath.Join(p.Root, "build.log"),
		filepath.Join(p.Root, ".flai-cache", "verify", "jobs", "S-0001-host-alias.log"): filepath.Join(home, ".flai", "other.log"),
	}
	for link, target := range links {
		if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, link); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Join(p.Root, ".flai-cache", "folder.log"), 0o755); err != nil {
		t.Fatal(err)
	}
	for name, given := range map[string]string{
		"a path outside":                       filepath.Join(outside, "secret.log"),
		"the project's own file":               "build.log",
		"a file not a .log":                    ".flai-cache/config.json",
		"flai's config":                        "~/.flai/config.json",
		"a .log beside flai's config":          "~/.flai/other.log",
		"a link out of .flai-cache":            ".flai-cache/out.log",
		"a .log link to flai's config":         ".flai-cache/config.log",
		"an agent run's name linked to config": "~/.flai/serve/agents/shipyard-S-0009-20261114T100000Z.log",
		"a folder linked out":                  ".flai-cache/kluster/secret.log",
		"a dangling link":                      ".flai-cache/dangling.log",
		"a link to another project's run":      "~/.flai/serve/agents/shipyard-S-0008-20261114T100000Z.log",
		"a link to the project's own file":     ".flai-cache/verify/jobs/S-0001-alias.log",
		"a link beside flai's config":          ".flai-cache/verify/jobs/S-0001-host-alias.log",
		"a .. out of the root":                 "../outside.log",
		"a .. out of .flai-cache":              ".flai-cache/../build.log",
		"a .. from ~/":                         "~/.flai/serve/../other.log",
		"a worktree":                           ".flai-cache/worktrees/S-0001/.flai-cache/x.log",
		"a folder linked to the worktrees":     ".flai-cache/merge-queue/S-0001/.flai-cache/x.log",
		"a dev flai's state in .flai-cache":    ".flai-cache/serve/agents/other-S-0001-20261114T100000Z.log",
		"another project's agent run":          "~/.flai/serve/agents/other-S-0001-20261114T100000Z.log",
		"a key with the project's as a prefix": "~/.flai/serve/agents/shipyard-ui-S-0001-20261114T100000Z.log",
		"another project's image job":          "~/.flai/serve/image-jobs/20261113T093000Z-001.log",
		"an image job no project has":          "~/.flai/serve/image-jobs/20261113T093000Z-002.log",
		"another project's checks run":         "~/.flai/serve/checks/S-0002-20261114T100000Z.log",
		"a checks run with no record":          "~/.flai/serve/checks/S-0003-20261114T100000Z.log",
		"a .log flai serve does not write":     "~/.flai/serve/journal.log",
		"a folder named .log":                  ".flai-cache/folder.log",
	} {
		got, e := readLogPart(t, host, p, given, 0)
		if e == nil {
			t.Errorf("%s: read %+v, want a refusal", name, got)
			continue
		}
		if e.Code != channel.CodeInvalidParams || !strings.Contains(e.Message, given+" cannot be read") || !strings.Contains(e.Message, "log.read reads only a .log file flai wrote for this project") {
			t.Errorf("%s: %d %s, want a refusal that names the path and says what is readable", name, e.Code, e.Message)
		}
		if strings.Contains(e.Message, "the content of") || strings.Contains(e.Message, home) || !strings.Contains(given, outside) && strings.Contains(e.Message, outside) {
			t.Errorf("%s: the refusal tells of the file: %s", name, e.Message)
		}
	}
}

func TestLogReadReadsNoServeLogWhereNoFlaiServeAnswers(t *testing.T) {
	p := shipyard(t)
	home, host := logsHome(t, nil)
	writeLog(t, filepath.Join(host.ServeDir, "serve.log"), "serve")
	writeLog(t, filepath.Join(home, ".flai", "host", "host.log"), "host")
	for _, given := range []string{"~/.flai/serve/serve.log", "~/.flai/host/host.log"} {
		if _, e := readLogPart(t, Host{}, p, given, 0); e == nil || e.Code != channel.CodeInvalidParams {
			t.Errorf("%s with no flai serve: %v, want a refusal", given, e)
		}
	}
	writeLog(t, filepath.Join(host.ServeDir, "image-jobs", "20261113T093000Z-001.log"), "job")
	if _, e := readLogPart(t, host, p, "~/.flai/serve/image-jobs/20261113T093000Z-001.log", 0); e == nil || !strings.Contains(e.Message, "runs no image jobs") {
		t.Errorf("an image job's log with no image jobs: %v, want a refusal that says so", e)
	}
}

func TestLogReadRefusesBadParams(t *testing.T) {
	p := shipyard(t)
	for params, want := range map[string]string{
		`{"offset":0}`:                     "give the path of the log",
		`{"path":"  ","offset":0}`:         "give the path of the log",
		`{"path":"a\nb.log","offset":0}`:   "line break or NUL",
		`{"path":"x.log","offset":-1}`:     "is not an offset",
		`{"path":"x.log","offset":"zero"}`: "",
	} {
		e := ask(t, Host{}, p, "log.read", params, nil)
		if e == nil || e.Code != channel.CodeInvalidParams || !strings.Contains(e.Message, want) {
			t.Errorf("%s: %v, want a bad-params refusal saying %q", params, e, want)
		}
	}
}
