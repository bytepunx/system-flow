---
id: T-1717
type: task
nature: feature
title: An integration test drives a story of a scratch repository to review on a stub opencode that replays the session, through flai mcp and flai guard
status: backlog
parent: S-0412
owner: alex
created: 2026-10-16T05:20:35Z
updated: 2026-11-07T09:34:08Z
transitions: []
stream: S-0412
tags: [cli]
touches: [flai/tests/integration/opencode_test.go]
after: [T-1716]
---
# T-1717 An integration test drives a story of a scratch repository to review on a stub opencode that replays the session, through flai mcp and flai guard

## Work

- Add `flai/tests/integration/opencode_test.go`, skipped under `testing.Short()` as the integration tier's tests are. Make a scratch system-flow repository with one ready story whose agent is `harness: opencode`, and set the host's program for `opencode` to the test binary re-run as a helper process, the stub.
- The stub reads the configuration the adapter wrote, replays `session.jsonl`: for each tool call it hands the plugin's input to `flai guard --harness opencode` and, when let through, makes the call through `flai mcp --agent <name>`; and it writes each event to standard output as `opencode run --format json` would.
- Start the agent through `flai serve`'s start path and wait for the stub to end. Assert the story was started, one task is done through `task_done`, and the story is in review.
- No `opencode` binary, network, or key is needed.
- Waits for T-1716: the session it replays.

## Done when

The test passes in the integration tier on a host with no `opencode` on `PATH`, and the story it drives ends in review.

## Notes
