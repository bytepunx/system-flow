package storygit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

// The conflict log (S-0337): one line for each conflict flai meets between
// stories, so that flai stats can count, per week, the conflicts found early
// by a trial merge apart from those met late by a rebase that stops. It is
// .flai-cache/conflicts.jsonl in the main checkout, beside overlaps.jsonl,
// append-only and never trimmed. The files flai resolves itself, generated
// and merged (ADR-0098, ADR-0126), are left out of each line's paths, and a
// stop on them alone is not a line.

// Where a conflict was met, as Conflict.Where says it.
const (
	// WhereTrial is a trial merge at a sync that sent or changed the conflict
	// message of a pair (ADR-0121).
	WhereTrial = "trial"
	// WhereSync is a rebase at flai stream sync, or flai task done, that
	// stopped on conflicts.
	WhereSync = "sync"
	// WhereAccept is the rebase of flai accept that stopped on conflicts.
	WhereAccept = "accept"
)

// Conflict is one line of the conflict log.
type Conflict struct {
	At    time.Time `json:"at"`              // when it was met, UTC
	Story string    `json:"story"`           // the story whose sync or acceptance met it
	Other string    `json:"other,omitempty"` // the other story of a trial merge's pair
	Where string    `json:"where"`           // WhereTrial, WhereSync, or WhereAccept
	Paths []string  `json:"paths"`           // the conflicting paths, sorted
}

// ConflictsPath is where the conflict log is in the cache folder cacheDir.
func ConflictsPath(cacheDir string) string { return filepath.Join(cacheDir, "conflicts.jsonl") }

// AppendConflict appends c to the conflict log in cacheDir as one line, its
// time in UTC to the second and its paths sorted.
func AppendConflict(cacheDir string, c Conflict) error {
	if c.Story == "" || !slices.Contains([]string{WhereTrial, WhereSync, WhereAccept}, c.Where) {
		return fmt.Errorf("log a conflict with story %q and where %q: a conflict needs a story and where it was met, one of %s, %s, %s", c.Story, c.Where, WhereTrial, WhereSync, WhereAccept)
	}
	c.At = c.At.UTC().Truncate(time.Second)
	c.Paths = append([]string{}, c.Paths...)
	slices.Sort(c.Paths)
	line, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("encode the conflict of %s for the conflict log: %w", c.Story, err)
	}
	path := ConflictsPath(cacheDir)
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return fmt.Errorf("make the cache folder %s for the conflict log: %w", cacheDir, err)
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open the conflict log %s: %w", path, err)
	}
	if _, err := f.Write(append(line, '\n')); err != nil {
		_ = f.Close()
		return fmt.Errorf("append to the conflict log %s: %w", path, err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("close the conflict log %s: %w", path, err)
	}
	return nil
}

// ReadConflicts returns every line of the conflict log in cacheDir, oldest
// first: none when there is no log yet, and an error when it cannot be read or
// a line cannot be parsed. Blank lines are skipped.
func ReadConflicts(cacheDir string) ([]Conflict, error) {
	path := ConflictsPath(cacheDir)
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read the conflict log %s: %w", path, err)
	}
	var out []Conflict
	for i, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var c Conflict
		if err := json.Unmarshal([]byte(line), &c); err != nil {
			return nil, fmt.Errorf("parse line %d of the conflict log %s: %w; correct or remove that line", i+1, path, err)
		}
		out = append(out, c)
	}
	return out, nil
}

// RecordConflict appends c to the conflict log in cacheDir with the files
// flai resolves itself left out of its paths, and appends nothing when no
// path is left. A line it cannot append is logged at warn: the sync or the
// acceptance that met the conflict goes on.
func RecordConflict(cacheDir string, c Conflict, files Files, log *slog.Logger) {
	c.Paths = withoutResolved(c.Paths, files)
	if len(c.Paths) == 0 {
		return
	}
	if err := AppendConflict(cacheDir, c); err != nil {
		if log == nil {
			log = slog.New(slog.DiscardHandler)
		}
		log.Warn("conflict not logged", "component", "storygit", "story", c.Story, "where", c.Where, "err", err)
	}
}

// recordConflict records c, met at o.Now, in the conflict log of o.Repo.
func recordConflict(o SyncOptions, c Conflict) {
	c.At = o.Now
	RecordConflict(o.Repo.CacheDir(), c, o.Files, o.Log)
}
