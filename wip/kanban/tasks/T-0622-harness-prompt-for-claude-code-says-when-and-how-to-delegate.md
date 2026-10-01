---
id: T-0622
type: task
nature: improvement
title: harness.Prompt for claude-code says when and how to delegate
status: done
parent: S-0175
owner: arobson
created: 2026-10-01T07:53:01Z
updated: 2026-10-01T08:09:06Z
transitions:
  - to: ready
    at: 2026-10-01T07:53:25Z
    by: agent-S-0175
  - to: in-progress
    at: 2026-10-01T08:08:28Z
    by: agent-S-0175
  - to: done
    at: 2026-10-01T08:09:06Z
    by: agent-S-0175
stream: S-0175
tags: []
touches: [flai/internal/harness]
usage:
  source: log
  seconds: 38
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 12
      output: 106
      cache_read: 1775990
      cache_write: 8814
      cost: 0.8005
---
# T-0622 harness.Prompt for claude-code says when and how to delegate

## Work

- The claude-code prompt says: hand code search, test and lint runs, and long logs to the explorer or verifier; before review have the verifier check the diff against the criteria and conventions; give a sub-agent the worktree, the story and task IDs, and the question; take back a summary; a question it returns is yours to ask with `thread_open`.
- Tests on the prompt.

## Done when

- `harness` tests pin the delegation paragraph for claude-code and its absence for the command harness.

## Notes
