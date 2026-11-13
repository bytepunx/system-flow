---
id: D-0592
title: flai stream sync --json prints a map that flai/cmd/stream.go builds by hand, not the JSON of storygit's sync result, so a new result field needs its key added there too
kind: decision
sources: [S-0702, S-0689, T-3030]
topics: [cli, go, flai, git]
paths: [flai/cmd/stream.go, flai/internal/storygit/sync.go]
date: 2026-11-11
status: active
absorbed_by: ""
asked_again: []
cited: [S-0751]
---

# D-0592 flai stream sync --json prints a map that flai/cmd/stream.go builds by hand, not the JSON of storygit's sync result, so a new result field needs its key added there too

`storygit`'s sync result carries JSON tags (`RegeneratedAtHead []string \`json:"regenerated_at_head"\`` in `flai/internal/storygit/sync.go`). But`flai stream sync --json` does not marshal that struct. `flai/cmd/stream.go` builds its own `map[string]any`, one for a stopped sync and one for a finished one, naming each key. A field added to the result alone never reaches the command's output.

S-0702 met this: "Widened T-3005's touches to flai/cmd/stream.go and flai/cmd/stream_sync.go: flai stream sync --json builds its own map, so the new regenerated_at_head key and a text line were added there, as docs/users/flai.md now says" (wip/archive/agents/S-0702.md, Decisions). The text report is `printSyncChecks` in the same package, which needs the new line too.

S-0689 met it again: "T-3030 widened to flai/cmd/stream.go and flai/internal/taskdone/taskdone.go: the stream sync --json output now carries rebase_error, abort_unsafe, and trial_merge_failed, and task_done's error for a StopGitError says what git said in place of 'conflicts in no path'" (wip/archive/agents/S-0689.md, Decisions). So a new kind of stop also needs its words in `flai/internal/taskdone/taskdone.go`.

A story that adds a sync result field puts `flai/cmd/stream.go` and `flai/cmd/stream_sync.go` in its touches, with `flai/internal/taskdone/taskdone.go` when it adds a stop, and documents the key in `docs/users/flai.md`.

Sources: S-0702, S-0689 (T-3030).
