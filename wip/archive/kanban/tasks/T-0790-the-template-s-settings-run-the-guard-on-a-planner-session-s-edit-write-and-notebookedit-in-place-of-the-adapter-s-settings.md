---
id: T-0790
type: task
nature: feature
title: The template's settings run the guard on a planner session's Edit, Write, and NotebookEdit, in place of the adapter's --settings
status: done
parent: S-0208
owner: alex
created: 2026-10-04T03:13:58Z
updated: 2026-10-04T03:34:55Z
transitions:
  - to: ready
    at: 2026-10-04T03:14:15Z
    by: agent-S-0208
  - to: in-progress
    at: 2026-10-04T03:14:29Z
    by: agent-S-0208
  - to: ready
    at: 2026-10-04T03:16:27Z
    by: agent-S-0208
  - to: in-progress
    at: 2026-10-04T03:34:55Z
    by: agent-S-0208
  - to: done
    at: 2026-10-04T03:34:55Z
    by: agent-S-0208
blocked:
  - from: 2026-10-04T03:16:28Z
    until: 2026-10-04T03:34:55Z
    reason: "the settings files are written by the operator: writes under .claude/ are refused to this session (TH-0096, I-0069)"
stream: S-0208
tags: []
touches: [".claude/settings.json", template/root/.claude/settings.json, flai/cmd/guard.go, flai/internal/harness]
after: [T-0786, T-0787]
usage:
  source: log
  seconds: 118
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 24
      output: 6867
      cache_read: 1718985
      cache_write: 29821
      cost: 0.662
---
# T-0790 The template's settings run the guard on a planner session's Edit, Write, and NotebookEdit, in place of the adapter's --settings

## Work

The operator's answer on TH-0096: update the settings files, so that the guard sees the planner's file edits. The planner writes epics, stories, tasks, their notes, and threads through flai only (`item_new`, `item_edit`, `thread_open`, `thread_reply`, and the CLI), as the operator's `planner.md` says; the guard T-0785 made already refuses its Edit, Write, and NotebookEdit everywhere.

- Both `.claude/settings.json` and `template/root/.claude/settings.json` gain a PreToolUse entry on `Edit|Write|NotebookEdit` that runs the guard only in a planner session (`FLAI_ROLE=plan`), so a story's agent pays no hook, and in this repository no rebuild of `bin/flai`, per edit. Writes under `.claude/` are refused to this session (I-0069), so the operator makes the change (TH-0096).
- The claude-code adapter stops passing a planner run `--settings`: the project's settings carry the hook.
- `flai guard --help` says the template's settings run it before the planner's file edits.

Waits for T-0787, which also changes `flai/cmd`.

## Done when

- [x] Both settings files run the guard on a planner session's Edit, Write, and NotebookEdit, and on nothing more for a story session
- [x] The adapter no longer passes `--settings`, and its test says so

## Notes
