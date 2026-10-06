---
id: T-1009
type: task
nature: improvement
title: This repository's .claude/settings.json runs the guard on sub-agent start and stop
status: done
parent: S-0285
owner: alex
created: 2026-10-06T10:38:25Z
updated: 2026-10-06T10:54:18Z
transitions:
  - to: ready
    at: 2026-10-06T10:54:17Z
    by: agent-S-0285
  - to: in-progress
    at: 2026-10-06T10:54:18Z
    by: agent-S-0285
  - to: done
    at: 2026-10-06T10:54:18Z
    by: agent-S-0285
stream: S-0285
tags: []
touches: [".claude/settings.json"]
after: [T-1004, T-1007]
usage:
  source: log
  seconds: 0
  models: []
---
# T-1009 This repository's .claude/settings.json runs the guard on sub-agent start and stop

## Work

This repository's `.claude/settings.json` runs `scripts/flai.sh guard` on `SubagentStart` and `SubagentStop`, as the template's runs `flai guard`. A write under `.claude/` needs the operator: the story's agent sends the whole file on a thread for them to paste, then tests it. Waits for T-1004 and T-1007, the guard and the template's file.

## Done when

- [ ] `.claude/settings.json` has the two hooks and is committed on the story branch

## Notes
