---
id: TH-0461
title: "I-0149's fix: one host upgrade, plus docs and a thread that names a published release"
anchor:
  path: wip/kanban/stories/S-0391-the-host-s-installed-flai-stays-behind-the-project-s-releases-after-i-0147-closed-so-the-mcp-task-done-still-refuses.md
  item: S-0391
status: resolved
participants: [agent-S-0391, orchestrator, alex]
created: 2026-10-15T05:08:59Z
updated: 2026-10-15T22:55:05Z
---

# TH-0461 I-0149's fix: one host upgrade, plus docs and a thread that names a published release

On wip/kanban/stories/S-0391-the-host-s-installed-flai-stays-behind-the-project-s-releases-after-i-0147-closed-so-the-mcp-task-done-still-refuses.md.

## Entries

### 2026-10-15T05:08:59Z agent-S-0391
**Cause.** The host runs flai 1.48.0 (`~/.flai/bin/flai`, `flai host` PID 26855, `flai serve` 26878). The fixes for this exist but shipped later: ADR-0158 (serve checks every minute, then upgrades or asks on a thread) in 1.50.1, and ADR-0160 (an outdated MCP `task_done` closes the task with the installed flai) in 1.51.1. A 1.48.0 host has neither, so nothing upgrades it and nothing asks you. No code change on main can reach a binary that is already running. flai 1.52.0 is published now.

**Recommendation.**

1. You run `flai host upgrade` once, in a shell. That clears the cause. ADR-0158 then keeps the host current: with `flai serve enable auto-upgrade` it upgrades itself, otherwise it asks you on one thread per release. I can't run it: the safety convention keeps me out of `~/.flai` and off your installed tools.
2. I build, on this branch:
   - a. Docs: a runbook step in `docs/operators/runbooks/update.md` and a row in `migrate.md` for hosts older than 1.50.1, which need this one manual upgrade. A fix to `docs/operators/index.md`, which says sessions started before an upgrade go on without a restart. That holds only for sessions whose MCP server is 1.51.1 or newer.
   - b. Code (ADR-0158 point 7): with `auto-upgrade` off, the upgrade thread names the newest *published* release, not the newest tag. A tag comes about 5 minutes before its binaries, so today the thread can name a `flai host upgrade --version X` that the host refuses. I'd record it as a new ADR refining ADR-0158, with a test.
3. After you upgrade, I check that the host and `flai serve` run ≥ 1.52.0 and that `flai task done` passes ADR-0142. Then I close I-0149, naming the upgrade and 2b.

I'm starting on 2a and 2b now. Reply "upgraded" once you've run it. Tell me if you want 2b left out, or if you'd rather turn on `auto-upgrade`.

### 2026-10-15T05:09:11Z orchestrator (recommendation)
**Recommendation:** run `flai host upgrade` once, and reply "upgraded" here. Keep 2b in S-0391. Choose `auto-upgrade` on or off as on TH-0460. Both threads ask for the same upgrade, so one answer settles both.

## Why keep 2b

- **It is the same cause.** With `auto-upgrade` off, ADR-0158's thread can name a tag whose binaries are not published yet, about five minutes early. Then `flai host upgrade --version X` refuses, and the host stays behind: that is I-0149 again, through the very path meant to prevent it.
- **It is small and tested.** One check, naming the newest published release, recorded as a new ADR refining ADR-0158, with a test.
- **It belongs to criterion 1,** "the cause I-0149 describes no longer occurs", better than to a separate story.

## Why this needs you

- The upgrade changes your machine. No agent may do it (`safety.md`).
- 2b adds code and an ADR beyond the planner's no-code plan on TH-0460. Its agent asked whether you want it.

Source: design/adrs/0158-flai-serve-checks-each-served-project-for-a-newer-flai-release-every-minute.md § Decision

### 2026-10-15T22:53:57Z alex
Confirmed the recommendation of 2026-10-15T05:09:11Z orchestrator.

Source: design/adrs/0158-flai-serve-checks-each-served-project-for-a-newer-flai-release-every-minute.md § Decision

### 2026-10-15T22:55:05Z agent-S-0391
Resolved: The operator upgraded the host to flai 1.52.0 and kept T-1599 (ADR-0161). The MCP task_done closed T-1601 with no refusal, and I-0149 is closed.
