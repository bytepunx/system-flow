---
id: T-0734
type: task
nature: feature
title: A story's review page offers a story for each open issue no open story links, by checkbox
status: done
parent: S-0198
owner: arobson
created: 2026-10-03T01:28:26Z
updated: 2026-10-03T02:09:26Z
transitions:
  - to: ready
    at: 2026-10-03T01:28:52Z
    by: agent-S-0198
  - to: in-progress
    at: 2026-10-03T01:58:25Z
    by: agent-S-0198
  - to: done
    at: 2026-10-03T02:09:26Z
    by: agent-S-0198
stream: S-0198
tags: []
touches: [flai/internal/hostapi, flaiover/src/routes/api/issues, flaiover/src/lib, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flai/cmd/issue.go, flai/cmd/issue_test.go, docs/users/flai-reference.md, docs/operators/settings.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-0733]
usage:
  source: log
  seconds: 661
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 150
      output: 846
      cache_read: 9378476
      cache_write: 191600
      cost: 3.9702
---
# T-0734 A story's review page offers a story for each open issue no open story links, by checkbox

## Work

- `flai/internal/hostapi` gets a read, `issue.list`, which runs `flai issue list --json` and passes `--story` through, and a write, `issue.story`, which runs `flai issue story <id>`. The contract test lists both.
- flaiover gets `GET /api/issues?story=S-nnnn` and `POST /api/issues/[id]/story` under `flaiover/src/routes/api/issues`.
- `Review.svelte` shows the open issues no open story links, each with a checkbox, as the designer chose on TH-0074, and a button that makes a backlog story for each checked one and says which stories it made. The logic that picks and orders the list lives in `flaiover/src/lib` with tests.
- Docs in the same change: `design/system/flaiover-dashboard.md` and `docs/users/dashboard.md`.

Waits for the fourth task, so that the command it runs and the docs it extends are settled.

## Done when

- hostapi tests cover both methods' arguments, and flaiover tests cover the list logic.
- `go test -race -short ./internal/hostapi/` and flaiover's tests for the changed files pass.

## Notes
