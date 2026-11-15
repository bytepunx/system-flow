---
id: T-3364
type: task
nature: feature
title: flai accept --resolve sends a story in review back to its agent to rebase and resolve what holds its acceptance, and the host method accept.resolve runs it
status: done
parent: S-0788
owner: alex
created: 2026-11-14T18:08:47Z
updated: 2026-11-14T18:34:27Z
transitions:
  - to: ready
    at: 2026-11-14T18:09:11Z
    by: agent-S-0788
  - to: in-progress
    at: 2026-11-14T18:09:11Z
    by: agent-S-0788
  - to: done
    at: 2026-11-14T18:34:27Z
    by: agent-S-0788
stream: S-0788
tags: []
touches: [flai/cmd/accept.go, flai/cmd/accept_resolve.go, flai/cmd/accept_resolve_test.go, flai/internal/hostapi/writes_items.go, flai/internal/hostapi/writes_items_test.go, docs/users/flai.md, docs/users/flai-reference.md, design/system/flai-cli.md, flaiover/src/lib/server/agent.ts, design/adrs/0302-the-operator-sends-a-review-story-that-cannot-be-accepted-back-to-its-agent-to.md, design/adrs/README.md, docs/operators/settings.md, flai/internal/hostapi/guardreads_test.go, flai/internal/hostapi/writes_test.go]
usage:
  source: log
  seconds: 1516
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 117
      output: 39417
      cache_read: 7046777
      cache_write: 363377
      cost: 4.373
---

# T-3364 flai accept --resolve sends a story in review back to its agent to rebase and resolve what holds its acceptance, and the host method accept.resolve runs it

## Work

Add `flai accept S-nnnn --resolve --by <operator>`. It accepts nothing. It runs the acceptance preview, as `flai accept --dry-run` does, and when the preview has blockers it sends the story back to its agent through the acceptance send-back of ADR-0248: it moves the story from review to in-progress as `storygit.SentBackAuthor` with a reason naming the blockers, and opens the send-back thread (title begins `storygit.SentBackPrefix`), or replies on the one still open. The thread lists each blocker as the preview words it and tells the agent to run `flai stream sync`, resolve what it stops on, take each blocker's recommended action, close out again, and move the story to review again. flai serve's `startAcceptSentBack` then starts the agent and `resolveAcceptSentBack` resolves the thread, with no change to flai serve. It refuses (exit 4) a story not in review, and one with no blockers, saying to accept it. `--json` answers the send-back object `sentBack` the acceptance already gives. Add the host method `accept.resolve`, `{id}`, which runs it with `--by` the project's owner. Record the decision as an ADR. Document the flag in `docs/users/flai.md` and `design/system/flai-cli.md`, and regenerate `docs/users/flai-reference.md`. No task waits before it: it is layer 1.

## Done when

- `flai accept S-nnnn --resolve` sends a review story with blockers back to in-progress with the thread, and refuses one with none or not in review, each with a test.
- `accept.resolve` builds that command line, with a test.
- `flai test` passes on the changed paths.

## Notes
