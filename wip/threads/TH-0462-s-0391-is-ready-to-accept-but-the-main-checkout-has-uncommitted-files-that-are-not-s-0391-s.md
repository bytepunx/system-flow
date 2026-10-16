---
id: TH-0462
title: S-0391 is ready to accept, but the main checkout has uncommitted files that are not S-0391's
anchor:
  path: wip/kanban/stories/S-0391-the-host-s-installed-flai-stays-behind-the-project-s-releases-after-i-0147-closed-so-the-mcp-task-done-still-refuses.md
  item: S-0391
status: resolved
participants: [orchestrator, alex]
created: 2026-10-15T23:26:05Z
updated: 2026-10-16T05:24:55Z
---

# TH-0462 S-0391 is ready to accept, but the main checkout has uncommitted files that are not S-0391's

On wip/kanban/stories/S-0391-the-host-s-installed-flai-stays-behind-the-project-s-releases-after-i-0147-closed-so-the-mcp-task-done-still-refuses.md.

## Entries

### 2026-10-15T23:26:05Z orchestrator
**Recommendation:** commit the files below on main yourself, or stash them. I then accept S-0391 at b5ca3495. I do not recommend `--yes`: it would fold the analyzer's report and your design-system files into S-0391's acceptance commit.

## The blocker

`flai accept S-0391 --dry-run` says "the real run refuses until they are committed or stashed, or --yes includes them". The files, uncommitted for 20 minutes or more:

- **The analyzer's run:**
  - `design/analysis/2026-10-15-bottlenecks.md` and `design/analysis/README.md`
  - `design/issues/I-0150` to `I-0153`, and `design/issues/summary.md`
- **A design system for this repository:** `design/system/design-system.md` and `design/system/design-system.tokens.json`. These look like `flai design init` output. I did not make them, and no open story claims them.

I may not commit or stash them myself.

## Where S-0391 stands

- `flai verify` passed every step at the head, b5ca3495.
- The verifier met both criteria:
  - 1: the upgrade thread waits for the published release (ADR-0161, with tests), and the operator docs cover hosts older than 1.50.1.
  - 2: I-0149 is closed.
- The narrative records your host upgrade: the host runs flai 1.52.0, and `task_done` closed T-1601 with no refusal.

## For the record

I-0150, which the analyzer has just filed, describes exactly this: "the orchestrator's acceptance stops on a dirty main".

### 2026-10-16T05:24:55Z alex
Resolved: S-0391 was accepted
