---
id: TH-1083
title: "The planner cannot start: image flai-agent-sf-claude-code:local is not on this host"
anchor:
  path: wip/kanban/epics/E-0034-an-operator-defines-the-board-s-states-transitions-rules-and-timers-in-a.md
  item: E-0034
status: open
participants: [orchestrator]
created: 2026-11-11T23:05:48Z
updated: 2026-11-11T23:05:48Z
---

# TH-1083 The planner cannot start: image flai-agent-sf-claude-code:local is not on this host

On wip/kanban/epics/E-0034-an-operator-defines-the-board-s-states-transitions-rules-and-timers-in-a.md.

## Entries

### 2026-11-11T23:05:48Z orchestrator
**Recommended:** build the image on the host, then I start E-0034's planner on my next run.

```sh
make agent-images
```

Or, over the harness image already built: `scripts/agent-image.sh --layer`.

The `plan` call for E-0034, first in `flai plan --candidates`, was refused with:

```text
could not be started: the image flai-agent-sf-claude-code:local, the project's layer for the inner harness claude-code that agent.image.layer in system-flow.yaml names, is not on this host (Error response from daemon: No such image: flai-agent-sf-claude-code:local)
```

The same refusal holds every planner run. Until the image is built, I start none for E-0034 or for backlog stories S-0659, S-0723, S-0725, and S-0726. Building the image changes the host, so it is yours to do.
