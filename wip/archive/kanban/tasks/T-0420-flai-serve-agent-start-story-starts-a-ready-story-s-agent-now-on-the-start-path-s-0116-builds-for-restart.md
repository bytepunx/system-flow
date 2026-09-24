---
id: T-0420
type: task
nature: feature
title: "flai serve agent start <story> starts a ready story's agent now, on the start path S-0116 builds for restart"
status: done
parent: S-0115
owner: alex
created: 2026-09-24T09:07:09Z
updated: 2026-09-24T09:19:33Z
transitions:
  - to: ready
    at: 2026-09-24T09:19:32Z
    by: agent-S-0115
  - to: in-progress
    at: 2026-09-24T09:19:33Z
    by: agent-S-0115
  - to: done
    at: 2026-09-24T09:19:33Z
    by: agent-S-0115
stream: S-0115
tags: []
touches: [flai/internal/serve, flai/cmd]
---
# T-0420 flai serve agent start <story> starts a ready story's agent now, on the start path S-0116 builds for restart

## Work

Build on the start path S-0116 exports for `flai serve agent restart` (TH-0010): story/S-0115 is based on story/S-0116 once its T-0417 lands. Add `flai serve agent start <story>`. It starts the agent the story names, or the host's command, as flai serve does when a story enters ready, records the run in `serve/agents.json`, and journals it, so the serving flai tracks it: the dot, the outcome, the start again on an answer. It starts whether or not the story was ready before serve started and whoever is attending. It refuses, saying why, when:
- the `agent` host action is off for the project;
- the story is not in ready;
- an agent is running for the story;
- nothing can start it (no harness and no command);
- the in-progress limit is full.

## Done when

- Behavior tests cover the start, a story ready before serve started, and each refusal.
- The started run is journalled like every agent start.
- `make test` and lint pass.

## Notes
