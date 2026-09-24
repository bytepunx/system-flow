---
id: T-0417
type: task
nature: remediation
title: "flai serve agent restart <story> starts a new agent for a story whose agent dropped or failed"
status: done
parent: S-0116
owner: alex
created: 2026-09-24T09:00:32Z
updated: 2026-09-24T09:10:01Z
transitions:
  - to: ready
    at: 2026-09-24T09:06:36Z
    by: agent-S-0116
  - to: in-progress
    at: 2026-09-24T09:06:36Z
    by: agent-S-0116
  - to: done
    at: 2026-09-24T09:10:01Z
    by: agent-S-0116
stream: S-0116
tags: []
touches: [flai/internal/serve, flai/cmd]
---
# T-0417 flai serve agent restart <story> starts a new agent for a story whose agent dropped or failed

## Work

Export a restart path from `flai/internal/serve` that uses the launcher's own start and records the run in `serve/agents.json`, so the serving flai tracks it: the dot, the outcome, the resume on an answer. A run no launcher waits for is settled by `settleOrphans`. Add `flai serve agent restart <story>`. It starts a new agent in a new session for a story in ready or in-progress whose newest run is not running and is not waiting on an answer. It refuses, saying why, when:
- the `agent` action is off for the project;
- the story is in another state;
- an agent is running for it;
- its agent waits on an answer;
- it had no agent from flai serve;
- nothing can start it;
- for a story in ready, the in-progress limit is full.

## Done when

- Behavior tests cover the start and each refusal.
- The started run is journalled like every agent start.
- `make test` and lint pass.

## Notes
