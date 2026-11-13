package kluster

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// writeStepLog writes text as the log of step in namespace ns under root.
func writeStepLog(t *testing.T, root, ns, step, text string) {
	t.Helper()
	path := LogPath(root, ns, step)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}

// ADR-0263: a step is build, build-<component>, load, deploy, or teardown,
// and nothing that could name another file.
func TestIsStep(t *testing.T) {
	for _, step := range []string{"build", "build-api", "build-a1-b2", "load", "deploy", "teardown"} {
		if !IsStep(step) {
			t.Errorf("%q is a step", step)
		}
	}
	for _, step := range []string{"", "images", "build-", "build-API", "build-a-", "../x", "build/../../kubeconfig", "deploy.log", "Deploy"} {
		if IsStep(step) {
			t.Errorf("%q is not a step", step)
		}
	}
}

// ADR-0263 §4: ReadLog answers the log's path from the root, its size, the
// text from the offset with the next offset, and the step's outcome in the
// namespace's record.
func TestReadLogFromAnOffset(t *testing.T) {
	root := t.TempDir()
	at := time.Date(2026, 11, 13, 9, 0, 0, 0, time.UTC)
	if err := writeRecord(root, Record{Story: "S-0001", Namespace: "s-0001", Outcome: OutcomeRunning, StartedAt: at, Steps: []Step{{Name: "build-api", Outcome: OutcomePassed}, {Name: StepLoad, Outcome: OutcomeRunning}}}); err != nil {
		t.Fatal(err)
	}
	writeStepLog(t, root, "s-0001", "build-api", "one\ntwo\n")
	got, err := ReadLog(root, "s-0001", "build-api", 0)
	if err != nil {
		t.Fatal(err)
	}
	want := LogPart{Namespace: "s-0001", Step: "build-api", Log: ".flai-cache/kluster/logs/s-0001/build-api.log", Exists: true, Offset: 0, Next: 8, Size: 8, Text: "one\ntwo\n", Outcome: OutcomePassed}
	if got != want {
		t.Errorf("from 0\n got %+v\nwant %+v", got, want)
	}
	if got, err = ReadLog(root, "s-0001", "build-api", 4); err != nil || got.Text != "two\n" || got.Offset != 4 || got.Next != 8 || got.Reset {
		t.Errorf("from the middle: %+v %v", got, err)
	}
	if got, err = ReadLog(root, "s-0001", "build-api", 8); err != nil || got.Text != "" || got.Next != 8 || got.Reset {
		t.Errorf("at the end: %+v %v", got, err)
	}
	writeStepLog(t, root, "s-0001", "deploy", "deploying\n")
	if got, err = ReadLog(root, "s-0001", "deploy", 0); err != nil || got.Outcome != "" || got.Text != "deploying\n" {
		t.Errorf("a step the record does not list has no outcome: %+v %v", got, err)
	}
}

// ADR-0263 §4: a read stops at 64 KB, before a character it would cut, and
// the next read from its next gives the rest.
func TestReadLogEndsOnAWholeCharacter(t *testing.T) {
	root := t.TempDir()
	text := strings.Repeat("a", LogChunk-1) + "é" + strings.Repeat("b", 10)
	writeStepLog(t, root, "s-0001", StepDeploy, text)
	first, err := ReadLog(root, "s-0001", StepDeploy, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Text) != LogChunk-1 || first.Next != LogChunk-1 || first.Size != int64(len(text)) {
		t.Fatalf("the first read stops before the character it would cut: len %d, next %d, size %d", len(first.Text), first.Next, first.Size)
	}
	rest, err := ReadLog(root, "s-0001", StepDeploy, first.Next)
	if err != nil {
		t.Fatal(err)
	}
	if first.Text+rest.Text != text || rest.Next != int64(len(text)) {
		t.Errorf("two reads give the whole log: next %d, rest %q", rest.Next, rest.Text)
	}
}

// ADR-0263 §4: a log shorter than the offset, or gone while the offset was
// above 0, is a reset, read again from 0.
func TestReadLogResets(t *testing.T) {
	root := t.TempDir()
	writeStepLog(t, root, "harbour", StepBuild, "new\n")
	got, err := ReadLog(root, "harbour", StepBuild, 100)
	if err != nil || !got.Reset || got.Offset != 0 || got.Text != "new\n" || got.Next != 4 {
		t.Errorf("a log shorter than the offset is read again from 0: %+v %v", got, err)
	}
	got, err = ReadLog(root, "harbour", StepLoad, 10)
	if err != nil || !got.Reset || got.Exists || got.Offset != 0 || got.Next != 0 {
		t.Errorf("a log gone while the offset was above 0: %+v %v", got, err)
	}
}

// ADR-0263 §4: a missing log answers exists false, empty, with its path.
func TestReadLogOfAMissingLog(t *testing.T) {
	got, err := ReadLog(t.TempDir(), "s-0009", StepTeardown, 0)
	want := LogPart{Namespace: "s-0009", Step: StepTeardown, Log: ".flai-cache/kluster/logs/s-0009/teardown.log"}
	if err != nil || got != want {
		t.Errorf("no log\n got %+v %v\nwant %+v", got, err, want)
	}
}

// A log flai cannot open is an error naming the step and the namespace.
func TestReadLogFails(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(LogPath(root, "s-0001", StepBuild), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadLog(root, "s-0001", StepBuild, 0); err == nil || !strings.Contains(err.Error(), "read the build log of the namespace s-0001") {
		t.Errorf("a log that is a folder: %v", err)
	}
}
