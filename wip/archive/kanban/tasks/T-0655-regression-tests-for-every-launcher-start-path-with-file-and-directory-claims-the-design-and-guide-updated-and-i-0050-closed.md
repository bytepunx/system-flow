---
id: T-0655
type: task
nature: remediation
title: Regression tests for every launcher start path with file and directory claims, the design and guide updated, and I-0050 closed
status: done
parent: S-0182
owner: arobson
created: 2026-10-01T09:16:25Z
updated: 2026-10-01T09:54:05Z
transitions:
  - to: ready
    at: 2026-10-01T09:16:34Z
    by: agent-S-0182
  - to: in-progress
    at: 2026-10-01T09:21:42Z
    by: agent-S-0182
  - to: done
    at: 2026-10-01T09:54:05Z
    by: agent-S-0182
stream: S-0182
tags: []
touches: [flai/internal/serve/hold_test.go, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0050-flai-serve-1-23-0-which-has-holds-started-an-agent-for-a-story-the-board-showed-held.md, design/issues/summary.md]
usage:
  source: log
  seconds: 491
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 176
      output: 49552
      cache_read: 14690331
      cache_write: 187353
      cost: 5.2741
---
# T-0655 Regression tests for every launcher start path with file and directory claims, the design and guide updated, and I-0050 closed

## Work

Add launcher tests through `look`: a ready story whose claim is a directory while an in-progress story claims a file in it, and the reverse, start nothing until the holder is accepted; the same with the holder in review and with a holder whose agent has started but which is still in ready. Describe in `design/system/workflow.md` that `resume` restarts only an open story and what an operator-started agent is told, and in `docs/users/flai.md` what a held story's Start agent does. Close I-0050 with the cause and the fix, or, if the designer's evidence says otherwise, record it there.

## Done when

The tests pass and fail with the hold check taken out of `look`; `make flai-test` passes; `flai check --strict` passes; I-0050 is closed or carries the note the story's Notes call for.

## Notes
