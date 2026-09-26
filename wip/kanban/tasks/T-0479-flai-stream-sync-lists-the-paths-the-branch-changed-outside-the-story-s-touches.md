---
id: T-0479
type: task
nature: feature
title: flai stream sync lists the paths the branch changed outside the story's touches
status: done
parent: S-0131
owner: alex
created: 2026-09-26T18:05:56Z
updated: 2026-09-26T18:10:50Z
transitions:
  - to: ready
    at: 2026-09-26T18:06:13Z
    by: agent-S-0131
  - to: in-progress
    at: 2026-09-26T18:09:52Z
    by: agent-S-0131
  - to: done
    at: 2026-09-26T18:10:50Z
    by: agent-S-0131
stream: S-0131
tags: []
touches: [flai/cmd/stream_sync.go]
---
# T-0479 flai stream sync lists the paths the branch changed outside the story's touches

## Work

- List the paths `story/<id>` changed since the main branch (merge-base diff) that no entry of the story's claim covers (touches of the story and its open tasks, components as their paths, as ADR-0046 reads them), leaving out the wip folder.
- Print them with the command to widen the claim (`flai touches <id> …`); `--json` carries `outside_touches`.
- An integration test with real git covers a change outside the claim and one inside it.

## Done when

- Sync names the paths outside the claim, and the test passes.

## Notes
