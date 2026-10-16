---
id: T-1601
type: task
nature: feature
title: Verify the host caught up after the operator's upgrade and close I-0149
status: done
parent: S-0391
owner: alex
created: 2026-10-15T05:09:40Z
updated: 2026-10-15T22:54:45Z
transitions:
  - to: ready
    at: 2026-10-15T05:09:59Z
    by: agent-S-0391
  - to: in-progress
    at: 2026-10-15T22:54:34Z
    by: agent-S-0391
  - to: done
    at: 2026-10-15T22:54:45Z
    by: agent-S-0391
stream: S-0391
tags: []
touches: [design/issues/I-0149-the-host-s-installed-flai-stays-behind-the-project-s-releases-after-i-0147-closed-so-the-mcp-task-done-still-refuses.md, design/issues/summary.md]
after: [T-1599, T-1600]
usage:
  source: log
  seconds: 11
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 10
      output: 2529
      cache_read: 622833
      cache_write: 39848
      cost: 0.4602
---

# T-1601 Verify the host caught up after the operator's upgrade and close I-0149

## Work

Once the operator has run `flai host upgrade` (TH-0461), check that `~/.flai/bin/flai version`, `flai host`, and `flai serve` run a flai that ADR-0142 does not refuse against the newest release in the project's history, and that the installed `flai task done` closes a task. Then close I-0149 with `flai issue close --commit`, naming the upgrade and T-1599. Waits for T-1599 and T-1600, so that the close names what the story built, and for the operator's upgrade.

## Done when

- The host's flai is at the project's newest release's minor version, observed with `flai version` and `pgrep`.
- I-0149 is closed with a reason naming the upgrade and T-1599.

## Notes
