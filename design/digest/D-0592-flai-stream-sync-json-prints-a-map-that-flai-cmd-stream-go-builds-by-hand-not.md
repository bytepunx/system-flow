---
id: D-0592
title: flai stream sync --json prints a map that flai/cmd/stream.go builds by hand, not the JSON of storygit's sync result, so a new result field needs its key added there too
kind: decision
sources: [S-0702]
topics: [cli, go, flai, git]
paths: [flai/cmd/stream.go, flai/internal/storygit/sync.go]
date: 2026-11-11
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0592 flai stream sync --json prints a map that flai/cmd/stream.go builds by hand, not the JSON of storygit's sync result, so a new result field needs its key added there too

`storygit`'s sync result carries JSON tags (`RegeneratedAtHead []string \`json:"regenerated_at_head"\`` in `flai/internal/storygit/sync.go`). But`flai stream sync --json` does not marshal that struct. `flai/cmd/stream.go` builds its own `map[string]any`, one for a stopped sync and one for a finished one, naming each key. A field added to the result alone never reaches the command's output.

S-0702 met this: "Widened T-3005's touches to flai/cmd/stream.go and flai/cmd/stream_sync.go: flai stream sync --json builds its own map, so the new regenerated_at_head key and a text line were added there, as docs/users/flai.md now says" (wip/archive/agents/S-0702.md, Decisions). The text report is `printSyncChecks` in the same package, which needs the new line too.

A story that adds a sync result field puts `flai/cmd/stream.go` and `flai/cmd/stream_sync.go` in its touches, and documents the key in `docs/users/flai.md`.

Source: S-0702.
