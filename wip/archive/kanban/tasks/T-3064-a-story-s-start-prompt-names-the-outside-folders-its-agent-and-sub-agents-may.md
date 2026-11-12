---
id: T-3064
type: task
nature: improvement
title: A story's start prompt names the outside folders its agent and sub-agents may reach, and says every other outside path is refused
status: done
parent: S-0701
owner: alex
created: 2026-11-10T10:35:42Z
updated: 2026-11-12T00:21:25Z
transitions:
  - to: ready
    at: 2026-11-12T00:05:46Z
    by: agent-S-0701
  - to: in-progress
    at: 2026-11-12T00:05:46Z
    by: agent-S-0701
  - to: done
    at: 2026-11-12T00:21:25Z
    by: agent-S-0701
stream: S-0701
tags: [harness]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go]
after: [T-3059]
usage:
  source: log
  seconds: 939
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 31
      output: 9261
      cache_read: 1773622
      cache_write: 79693
      cost: 1.0015
---
# T-3064 A story's start prompt names the outside folders its agent and sub-agents may reach, and says every other outside path is refused

## Work

In `flai/internal/harness/harness.go`, add one paragraph to the story's start prompt, after its worktree lines. With outside folders in the start request, it names them and says the agent and its sub-agents may read and write there. In either case, it says that Read, Grep, Glob, Write, and Monitor are refused on any other path outside the project. Such work then uses Bash, and keeps its files under the worktree's git-ignored `bin/<story>/`, as D-0480 records. In `delegation`, tell the agent to name the folders to a sub-agent beside the worktree's path.

This task waits for T-3059, which adds the folders to the start request. It runs beside the container task, with which it shares no file.

## Done when

- Tests in `flai/internal/harness/harness_test.go` show the paragraph with two folders, and the refusal line with none.
- `flai test` passes on the changed paths.

## Notes
