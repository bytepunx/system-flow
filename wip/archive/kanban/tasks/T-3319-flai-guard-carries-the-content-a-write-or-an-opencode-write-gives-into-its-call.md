---
id: T-3319
type: task
nature: improvement
title: flai guard carries the content a Write or an OpenCode write gives into its Call
status: done
parent: S-0730
owner: alex
created: 2026-11-13T19:53:57Z
updated: 2026-11-13T20:07:07Z
transitions:
  - to: ready
    at: 2026-11-13T19:59:21Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T19:59:21Z
    by: system-flow
  - to: done
    at: 2026-11-13T20:07:07Z
    by: agent-S-0730
stream: S-0730
tags: [flai]
touches: [flai/internal/guard/call.go, flai/internal/guard/claudecode.go, flai/internal/guard/claudecode_test.go, flai/internal/guard/opencode.go, flai/internal/guard/opencode_test.go]
usage:
  source: log
  seconds: 466
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 8
      output: 1098
      cache_read: 921673
      cache_write: 8825
      cost: 0.2723
---
# T-3319 flai guard carries the content a Write or an OpenCode write gives into its Call

## Work

The guard drops the content of a whole-file write today. `Input` in `flai/internal/guard/claudecode.go` keeps only `command`, `file_path`, `notebook_path`, and flai's args. `openCodeArgs` in `flai/internal/guard/opencode.go` keeps only `command`, `filePath`, and `patchText`. So the rule that I-0352 needs cannot see what a Write would leave in the file.

- Add a `Content` field to `Call` in `flai/internal/guard/call.go`: the whole text a FileWrite leaves in `Path`, set only for a tool that replaces the file whole.
- Read Claude Code's `Write` `tool_input.content` into it in `claudecode.go`. Leave `Edit`, `MultiEdit`, and `NotebookEdit` without content: they change part of a file.
- Read OpenCode's `write` `content` into it in `opencode.go`. Leave `edit`, `multiedit`, and `apply_patch` without it.
- Decide nothing new here. The guard still fails open on input it cannot read.

Waits for no task: it is the first layer.

## Done when

- `claudecode_test.go` shows a `Write` hook input's content reaching `Call.Content`, and an `Edit` input leaving it empty.
- `opencode_test.go` shows the same for OpenCode's `write` and `edit`.
- `flai test` on the changed paths passes.

## Notes
