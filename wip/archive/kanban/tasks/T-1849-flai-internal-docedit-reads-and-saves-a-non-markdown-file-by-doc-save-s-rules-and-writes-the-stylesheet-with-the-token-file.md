---
id: T-1849
type: task
nature: feature
title: flai/internal/docedit reads and saves a non-markdown file by doc.save's rules, and writes the stylesheet with the token file
status: done
parent: S-0436
owner: alex
created: 2026-10-16T21:31:30Z
updated: 2026-10-16T22:05:47Z
transitions:
  - to: ready
    at: 2026-10-16T21:53:26Z
    by: agent-S-0436
  - to: in-progress
    at: 2026-10-16T21:53:27Z
    by: agent-S-0436
  - to: done
    at: 2026-10-16T22:05:47Z
    by: agent-S-0436
stream: S-0436
tags: [cli]
touches: [flai/internal/docedit/file.go, flai/internal/docedit/file_test.go, flai/internal/docedit/docedit.go]
usage:
  source: log
  seconds: 740
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 45
      output: 17612
      cache_read: 2500101
      cache_write: 82425
      cost: 1.3528
---
# T-1849 flai/internal/docedit reads and saves a non-markdown file by doc.save's rules, and writes the stylesheet with the token file

## Work

Add `file.go` to `flai/internal/docedit`, beside `docedit.go`, reusing its hash, check, restore, and commit rules:

- `FileReach(m, path)`: a file under the manifest's design or docs folder, or `m.DesignTokens()` or `m.DesignStylesheet()`. Refuse the wip folder, a `.md` file, a binary file (a NUL byte in the first 8 KB), and a file over the size limit, each with its reason. Recommended limit: 1 MiB, the largest token file plus room; record the figure in the story's `## Decisions`.
- `FileShow` and `FileSave`, each with an optional story worktree as ADR-0145 has it.
- When the path is the token file and the manifest names a stylesheet, `FileSave` writes the stylesheet from the new tokens with `designsystem.Stylesheet`, checks with both in place, and commits both (ADR-0150).

Waits for nothing: it is the rule every other task calls.

## Done when

Tests cover each reach and each refusal, a stale hash, a finding that restores both files, a commit with both paths, a save with autocommit off, and a story's worktree. `flai test flai/internal/docedit` passes.

## Notes

`docedit.go` is read, not changed; if a shared helper must move, widen the touches with it.
