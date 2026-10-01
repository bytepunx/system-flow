---
id: T-0675
type: task
nature: improvement
title: A verifier never makes the corrections it finds
status: done
parent: S-0189
owner: arobson
created: 2026-10-01T10:53:32Z
updated: 2026-10-01T10:54:11Z
transitions:
  - to: ready
    at: 2026-10-01T10:53:48Z
    by: agent-S-0189
  - to: in-progress
    at: 2026-10-01T10:53:49Z
    by: agent-S-0189
  - to: done
    at: 2026-10-01T10:54:11Z
    by: agent-S-0189
stream: S-0189
tags: []
touches: [flai/internal/guard]
usage:
  source: log
  seconds: 22
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 34
      output: 169
      cache_read: 1008948
      cache_write: 18480
      cost: 0.4769
---
# T-0675 A verifier never makes the corrections it finds

## Work

Criterion 3. The wording is already committed with T-0672. Add a guard test that a verifier's commands that would correct the worktree or record a fix are refused: `git add`, `git commit`, `git apply`, `git checkout --`, `git restore`, and `flai move`.

## Done when

- The test passes.

## Notes
