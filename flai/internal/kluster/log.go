package kluster

// A step's log, read a part at a time from an offset (ADR-0263 §4): the
// host method kluster.log answers one part for the dashboard's live feed,
// and flai kluster log reads every part from the start (S-0763).

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"unicode/utf8"
)

// LogChunk is the most of a step's log one ReadLog answers.
const LogChunk = 64 << 10

// stepName is a deploy's, a promotion's, or a teardown's step: build, load,
// deploy, teardown, or build-<component> (ADR-0263).
var stepName = regexp.MustCompile(`^(build|load|deploy|teardown|build-[a-z0-9]([-a-z0-9]{0,61}[a-z0-9])?)$`)

// IsStep says whether step names a step whose log ReadLog reads.
func IsStep(step string) bool { return stepName.MatchString(step) }

// LogPart is a part of a step's log, read from Offset, with where to read
// from next: what the host method kluster.log answers.
type LogPart struct {
	Namespace string `json:"namespace"`
	Step      string `json:"step"`
	// Log is the log's path from the project's root.
	Log string `json:"log"`
	// Exists says whether the log is there; with none, Text is empty.
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
	// Reset is true when the log was shorter than the offset asked, as when
	// a new job cleared it, and Text then starts at 0.
	Reset bool `json:"reset,omitempty"`
	// Outcome is the step's outcome in the namespace's record, empty when the
	// record has no such step.
	Outcome string `json:"outcome,omitempty"`
}

// ReadLog reads at most LogChunk bytes of namespace ns's step log under
// root from offset, ending on a whole character, with the step's outcome
// from the namespace's record; ns is a DNS label, step is one IsStep
// accepts, and offset is 0 or more, as the caller checks.
func ReadLog(root, ns, step string, offset int64) (LogPart, error) {
	path := LogPath(root, ns, step)
	rel, err := filepath.Rel(root, path)
	if err != nil {
		rel = path
	}
	out := LogPart{Namespace: ns, Step: step, Log: filepath.ToSlash(rel), Offset: offset, Next: offset}
	if r, ok, err := ReadRecord(root, ns); err == nil && ok {
		for _, st := range r.Steps {
			if st.Name == step {
				out.Outcome = st.Outcome
			}
		}
	}
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		out.Offset, out.Next = 0, 0
		out.Reset = offset > 0
		return out, nil
	}
	if err != nil {
		return LogPart{}, fmt.Errorf("read the %s log of the namespace %s: %w", step, ns, err)
	}
	defer func() { _ = f.Close() }()
	st, err := f.Stat()
	if err != nil {
		return LogPart{}, fmt.Errorf("read the %s log of the namespace %s: %w", step, ns, err)
	}
	out.Exists, out.Size = true, st.Size()
	if offset > out.Size {
		out.Offset, out.Reset = 0, true
	}
	buf := make([]byte, min(out.Size-out.Offset, LogChunk))
	n, err := f.ReadAt(buf, out.Offset)
	if err != nil && !errors.Is(err, io.EOF) {
		return LogPart{}, fmt.Errorf("read the %s log of the namespace %s: %w", step, ns, err)
	}
	buf = wholeRunes(buf[:n])
	out.Text, out.Next = string(buf), out.Offset+int64(len(buf))
	return out, nil
}

// wholeRunes is b without a character cut short at its end, so the next
// read starts where that character does.
func wholeRunes(b []byte) []byte {
	for i := 1; i <= utf8.UTFMax && i <= len(b); i++ {
		if utf8.RuneStart(b[len(b)-i]) {
			if !utf8.FullRune(b[len(b)-i:]) {
				return b[:len(b)-i]
			}
			return b
		}
	}
	return b
}
