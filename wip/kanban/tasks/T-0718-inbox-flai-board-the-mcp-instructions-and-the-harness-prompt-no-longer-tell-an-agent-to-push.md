---
id: T-0718
type: task
nature: improvement
title: Inbox, flai board, the MCP instructions, and the harness prompt no longer tell an agent to push
status: done
parent: S-0195
owner: arobson
created: 2026-10-02T23:31:03Z
updated: 2026-10-02T23:46:36Z
transitions:
  - to: ready
    at: 2026-10-02T23:31:26Z
    by: agent-S-0195
  - to: in-progress
    at: 2026-10-02T23:40:30Z
    by: agent-S-0195
  - to: done
    at: 2026-10-02T23:46:36Z
    by: agent-S-0195
stream: S-0195
tags: []
touches: [flai/internal/mcpserver, flai/internal/workitem/boardview.go, flai/cmd/board.go, flai/internal/harness]
after: [T-0717]
usage:
  source: log
  seconds: 366
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 126
      output: 984
      cache_read: 5994912
      cache_write: 135434
      cost: 2.5152
---
# T-0718 Inbox, flai board, the MCP instructions, and the harness prompt no longer tell an agent to push

## Work

Remove the push duty and keep the facts (ADR-0067 as accepted): the single-project and folder MCP servers' instructions (`server.go`, `folder.go`) stop telling agents to run `flai push --pending`, and say instead that publishing is the operator's, done by an agent only when the operator asks. `inbox` and `flai board` (text and JSON, `BoardView`) stop carrying `unpushed` with its push command, and instead report what is accepted and not yet published (`unpublished`: the accepted items pending a publish, from `release.PendingIDs`), described as information for the operator or an agent the operator asks to publish, not a duty. Check that `harness.Prompt` says nothing of pushing, and pin it in a test. Rewrite `unpushed_test.go` to assert the instructions and inbox carry no push duty and do report what is unpublished. Waits for T-0717: it builds what ADR-0067 decides.

## Done when

- The instructions, `inbox`, `flai board`, and `harness.Prompt` say nothing that tells an agent to push an unpushed acceptance
- `inbox` and `flai board` still say which accepted items are not yet published
- Tests pin both

## Notes
