package guard

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
)

// The Claude Code hook events on a sub-agent, which the guard tells from a
// PreToolUse by hook_event_name. An input that names neither is a tool call,
// whether it names PreToolUse or, as the guard was first fed, no event.
const (
	HookSubagentStart = "SubagentStart"
	HookSubagentStop  = "SubagentStop"
)

// RecordDir is the folder, under a project's .flai-cache in its main
// checkout, that holds each session's record of its running sub-agents.
const RecordDir = "guard"

// SubAgent is one sub-agent running in a session, as its SubagentStart
// reported it: its ID, the name of its definition, and when it started.
type SubAgent struct {
	ID      string    `json:"agent_id"`
	Type    string    `json:"agent_type"`
	Started time.Time `json:"started"`
}

// record is what a session's record file holds.
type record struct {
	Running []SubAgent `json:"running"`
}

// sessionName is what a session ID must be to name a record file.
var sessionName = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

// recordLockWait is how long Record waits for another flai guard to finish
// with a session's record; recordLockAge is how old a lock is when it is
// taken as left behind by a flai guard that died holding it. A record is held
// for a read and a write of a few lines, so either is far longer than any
// holder needs.
const (
	recordLockWait = 10 * time.Second
	recordLockAge  = 5 * time.Second
)

// Record notes in dir, the project's .flai-cache/guard, that e's sub-agent
// started (SubagentStart) or stopped (SubagentStop) at at, in the record of
// e's session. Each sub-agent of a layer starts in a hook process of its
// own, often in the same second, so the record is read and written under a
// lock and replaced in one step. A record left with no sub-agent running is
// removed. An event of another kind, a session ID that is not a safe file
// name, or an event with no agent ID is ignored.
func Record(dir string, e Event, at time.Time) error {
	if (e.HookEventName != HookSubagentStart && e.HookEventName != HookSubagentStop) || !sessionName.MatchString(e.SessionID) || e.AgentID == "" {
		return nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("record of session %s: %w", e.SessionID, err)
	}
	path := filepath.Join(dir, e.SessionID+".json")
	unlock, err := lockRecord(path + ".lock")
	if err != nil {
		return fmt.Errorf("record of session %s: %w", e.SessionID, err)
	}
	defer unlock()
	r, err := readRecord(path)
	if err != nil {
		return err
	}
	r.Running = slices.DeleteFunc(r.Running, func(s SubAgent) bool { return s.ID == e.AgentID })
	if e.HookEventName == HookSubagentStart {
		r.Running = append(r.Running, SubAgent{ID: e.AgentID, Type: e.AgentType, Started: at.UTC().Truncate(time.Second)})
	}
	if len(r.Running) == 0 {
		if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("record of session %s: %w", e.SessionID, err)
		}
		return nil
	}
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	if err := atomicfile.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("record of session %s: %w", e.SessionID, err)
	}
	return nil
}

// Running is the sub-agents the record in dir lists as running in session,
// in the order they started; none when the session has no record or its ID
// is not a safe file name.
func Running(dir, session string) ([]SubAgent, error) {
	if !sessionName.MatchString(session) {
		return nil, nil
	}
	r, err := readRecord(filepath.Join(dir, session+".json"))
	return r.Running, err
}

// readRecord reads a session's record; a missing one lists no sub-agent.
func readRecord(path string) (record, error) {
	var r record
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return r, nil
	}
	if err != nil {
		return r, err
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return r, fmt.Errorf("%s: %w: remove it, and the guard starts a new record at the session's next sub-agent", path, err)
	}
	return r, nil
}

// lockRecord takes the lock at path, a file one writer at a time creates,
// waiting while another holds it, and taking one older than recordLockAge.
func lockRecord(path string) (unlock func(), err error) {
	deadline := time.Now().Add(recordLockWait)
	for {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if err == nil {
			if err := f.Close(); err != nil {
				_ = os.Remove(path)
				return nil, err
			}
			return func() { _ = os.Remove(path) }, nil
		}
		if !errors.Is(err, os.ErrExist) {
			return nil, err
		}
		if st, err := os.Stat(path); err == nil && time.Since(st.ModTime()) > recordLockAge {
			_ = os.Remove(path)
			continue
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("%s is held by another flai guard after %s: remove it if none runs", path, recordLockWait)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
