---
id: T-2878
type: task
nature: feature
title: A story's prime pack in a kluster project names the deploy and checks commands
status: done
parent: S-0662
owner: alex
created: 2026-11-10T08:26:53Z
updated: 2026-11-13T13:56:21Z
transitions:
  - to: ready
    at: 2026-11-13T13:44:02Z
    by: agent-S-0662
  - to: in-progress
    at: 2026-11-13T13:44:02Z
    by: agent-S-0662
  - to: done
    at: 2026-11-13T13:56:21Z
    by: agent-S-0662
stream: S-0662
tags: [cli]
touches: [flai/internal/context/story.go, flai/internal/context/context.go, flai/internal/context/context_test.go]
usage:
  source: log
  seconds: 739
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 31
      output: 12030
      cache_read: 2311463
      cache_write: 98985
      cost: 1.2434
---
# T-2878 A story's prime pack in a kluster project names the deploy and checks commands

## Work

In `ForStory` (`flai/internal/context/story.go`), read the manifest's `kluster:` block. When it names a `deploy`, give the pack a field that names the commands a story's agent deploys and checks its namespace with, the CLI command and the MCP tool S-0660 and S-0661 added, and the namespace. Print it in `Pack.Header` (`context.go`) and in the JSON the MCP `prime` returns.

A project with no block, or no `deploy`, gets no such field or line. It reads the manifest only, not the deploy record, so it waits for no task.

## Done when

- A test in `context_test.go` primes a story in a project whose block names a `deploy` and finds the commands and the namespace in the header and the JSON.
- A test primes one with no block, and one with no `deploy`, and finds neither.
- `flai test` passes on the three files.

## Notes

- Criterion 7 of S-0662. Nothing in the pack reads the manifest today.
