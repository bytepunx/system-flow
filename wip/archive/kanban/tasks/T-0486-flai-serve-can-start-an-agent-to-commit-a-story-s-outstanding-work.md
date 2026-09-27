---
id: T-0486
type: task
nature: remediation
title: flai serve can start an agent to commit a story's outstanding work
status: done
parent: S-0140
owner: alex
created: 2026-09-26T21:11:25Z
updated: 2026-09-26T21:19:32Z
transitions:
  - to: ready
    at: 2026-09-26T21:11:31Z
    by: agent-S-0140
  - to: in-progress
    at: 2026-09-26T21:15:15Z
    by: agent-S-0140
  - to: done
    at: 2026-09-26T21:19:32Z
    by: agent-S-0140
stream: S-0140
tags: []
touches: [flai/internal/serve, flai/internal/hostapi, flai/internal/harness, flai/cmd]
---
# T-0486 flai serve can start an agent to commit a story's outstanding work

## Work

- A host method and `flai serve agent commit <story>` start the story's agent with a prompt to commit everything outstanding in its worktree and change nothing else, for a story in review or in progress whose worktree has uncommitted changes and whose agent is not running.
- The run is recorded and journalled like a restart.

## Done when

- [x] Behaviour tests cover the start, its prompt, and each refusal.

## Notes
