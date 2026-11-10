---
id: TH-0997
title: "Plan for S-0667: find why a hand docker run mounts the checkout empty, then fix and check it in scripts/agent-container.sh, then document it and close I-0276"
anchor:
  path: wip/kanban/stories/S-0667-a-docker-run-from-a-story-s-agent-s-shell-mounts-the-main-checkout-as-an-empty.md
  item: S-0667
status: resolved
participants: [planner-S-0667, orchestrator]
created: 2026-11-10T08:49:39Z
updated: 2026-11-10T08:54:04Z
---

# TH-0997 Plan for S-0667: find why a hand docker run mounts the checkout empty, then fix and check it in scripts/agent-container.sh, then document it and close I-0276

On wip/kanban/stories/S-0667-a-docker-run-from-a-story-s-agent-s-shell-mounts-the-main-checkout-as-an-empty.md.

## Entries

### 2026-11-10T08:49:39Z planner-S-0667
Plan for S-0667. Three tasks in three layers, one after another.

| Layer | Task | Nature | Touches | After |
|-------|------|--------|---------|-------|
| 1 | T-2926: find the cause of I-0276's empty mount on the host and propose a fix | research | the I-0276 issue file | none |
| 2 | T-2927: a hand-run container sees the checkout's files, and refuses with the cause when a mount is empty | remediation | `scripts/agent-container.sh` | T-2926 |
| 3 | T-2929: the runbook and `design/tech/docker.md` give the cause, and I-0276 is closed | improvement | `docs/operators/runbooks/agent-container.md`, `design/tech/docker.md`, the I-0276 file, `design/issues/summary.md` | T-2927 |

Figures:

- Forecast: 1h30m, raised from flai's 13m for the live docker runs T-2926 needs.
- Cost of delay: 50 USD a week, flai's figure from the 20m input.
- No `after` on other stories.

Assumptions:

- The story's agent runs on the host with `docker` on PATH. A planner or agent inside the container has no Docker socket, so it cannot reproduce the issue.
- The cause is not known yet. S-0567's agent ran on the host, and its hand bind mount came up empty inside and outside Claude Code's sandbox. `flai serve`'s container mounted the same path with its files. My guess is a difference in docker context or `DOCKER_HOST` between the two, whose daemon cannot see the WSL path. T-2926 tests this before anything is built.
- If T-2926 finds the fix is outside the repository, such as a Docker Desktop setting, it asks on a thread on this story first. T-2927 then adds only the check that refuses an empty mount and names the cause.
- "A test that reproduces it where one fits": a live docker reproduction cannot run in the test tiers. T-2927's check reproduces it on the host instead, and ShellCheck covers the script.

Nothing here needs your answer before the story is pulled.

### 2026-11-10T08:53:39Z orchestrator
Approved as written: T-2926, then T-2927, then T-2929; 1h30m, 50 USD a week.

- The cause is not known, so a research task on the host comes first and the fix builds on what it finds.
- T-2927's check refuses an empty mount and names the cause. Where a live Docker run cannot sit in the test tiers, that check is the reproduction criterion 1 allows.
- Asking on a thread before T-2927 when the fix lies outside the repository is right: a host Docker setting is yours.
- The five touches are files, each layer one task. No `after` on S-0670 or S-0682 is right: the trial merge judges the shared script.

### 2026-11-10T08:54:04Z orchestrator
Resolved: Plan approved as written by the orchestrator under plan_backlog_stories.
