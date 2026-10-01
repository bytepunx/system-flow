---
id: T-0664
type: task
nature: feature
title: flai serve agent restart starts a story with no run on this host, and its agent is told the story was begun elsewhere
status: in-progress
parent: S-0177
owner: alex
created: 2026-10-01T10:06:57Z
updated: 2026-10-01T10:11:31Z
transitions:
  - to: ready
    at: 2026-10-01T10:07:15Z
    by: agent-S-0177
  - to: in-progress
    at: 2026-10-01T10:11:31Z
    by: agent-S-0177
stream: S-0177
tags: []
touches: [flai/internal/serve, flai/internal/harness, flai/cmd/serve_actions.go]
---
# T-0664 flai serve agent restart starts a story with no run on this host, and its agent is told the story was begun elsewhere

## Work

- `serve.Restart` no longer refuses a story in ready or in progress that has no run here. A ready story starts, or queues when the limit is full or a hold applies. A story in progress starts with `harness.Request.Begun`. It still refuses while this host's agent runs or waits for an answer.
- `Begun` comes from the story's last move to in-progress (by, at), the narrative's agent and host, and the story's threads answered since then. `harness.Prompt` tells the agent where, when, and by whom the story was begun. It tells it to reconcile through `flai stream open`, read the narrative and tasks, go on from what is committed, and read the answered threads.
- Tests: restart of an in-progress story with no run, the prompt's text, and the refusal while this host's agent runs.
- `design/system/flai-cli.md`, the command's help, and `docs/users/flai.md` say so.

## Done when

The tests pass and the change is committed.
