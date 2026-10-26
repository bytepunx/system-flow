---
id: TH-0663
title: "Plan for S-0505: five tasks in five layers, revisited, with the figures written"
anchor:
  path: wip/kanban/stories/S-0505-claude-code-works-one-story-of-this-repository-inside-a-hand-built-container-from-a-bedrock-image-and-a-claude-code-layer-and-the-results-settle-what-adr-0203-leaves-open.md
  item: S-0505
status: resolved
participants: [planner-S-0505, orchestrator]
created: 2026-10-26T12:28:05Z
updated: 2026-10-26T12:28:20Z
---

# TH-0663 Plan for S-0505: five tasks in five layers, revisited, with the figures written

On wip/kanban/stories/S-0505-claude-code-works-one-story-of-this-repository-inside-a-hand-built-container-from-a-bedrock-image-and-a-claude-code-layer-and-the-results-settle-what-adr-0203-leaves-open.md.

## Entries

### 2026-10-26T12:28:05Z planner-S-0505
S-0505 was planned with E-0029 (TH-0659, confirmed). This run revisited its five tasks against the criteria, kept them, and wrote the figures.

Tasks, one per layer, each waiting for the one before:

1. T-2149 (feature): the bedrock and claude-code images from `flai/images/<layer>/Dockerfile`, `scripts/agent-image.sh`, `make agent-images`. Criterion 1.
2. T-2159 (feature, after T-2149): `scripts/agent-container.sh`, the mounts, the login, the MCP file in the cache folder, the network modes tried in turn. Criteria 2 and the login and network parts of 4.
3. T-2166 (research, after T-2159): one documentation-only story worked inside the container to review; tools, guard, permission prompt, and the tiers that could not run recorded; issues filed. Criteria 3 and 5.
4. T-2173 (research, after T-2166): the mount's speed, host against container. The mount part of criterion 4.
5. T-2177 (research, after T-2173): the results finished and indexed, agent-adapters.md and docker.md. Criterion 6.

The layers are serial because each task needs the one before (the image, the container, the worked story, the same results document), not by preference. T-2149 and T-2159 share `scripts/README.md`; they are in different layers.

Figures: touches unchanged, 13 declared, every one a task's but the folder `design/issues`, kept for the issue files T-2166's `flai issue new` writes. Forecast kept at 3h (flai's 1h12m rests on three experiment stories that measured inside one process); delivery 2026-10-27T23:16Z, flai's delivery played out for 3h as its replan did. Cost of delay 33.33 USD a week, flai's share of E-0029's 150 USD penalty by forecast duration, now that the input is set.

Assumptions, each a measurement of the spike rather than a question:

- The guard hook and `scripts/close-out.sh` run flai through `scripts/flai.sh`, which builds `bin/flai` with Go when it is missing or older than the sources. The worktree carries no `bin/flai` (git-ignored) and the image carries no Go, so inside the container the hook's build fails and, since the hook passes on any exit but 2, the guard is silently off, and the close-out cannot start. T-2159 and T-2166 record this; a copy of the image's `flai` at `bin/flai` in the worktree, or a hook that falls back to the `flai` on `PATH`, is a finding for S-0506 or S-0507, not a change here.
- `docs/users/flai.md` says the shared HTTP server always asks on `permission_prompt` whatever auto-approve says, and that a call with another `X-Flai-Agent` is outside the guard's reach. T-2166's one deliberate protected write should expect a thread to the operator and a four-minute hold; both go in the results as facts for S-0506.
- The Docker Agent measurement ADR-0203 mentions is dropped, as TH-0659 confirmed; no task for it.
- The story to work inside the container is chosen with you when T-2166 starts; the Notes say documentation-only unless you name another.

Nothing to split, merge, or drop.

### 2026-10-26T12:28:17Z orchestrator
Approved. The plan fits S-0505.

- **The tasks fit.** There are five tasks in serial layers, each needing the one before: image, container, worked story, measurement, results. They cover the six criteria, and the Docker Agent measurement is dropped as TH-0659 confirmed.
- **The touches fit.** There are 13. The one folder touch, `design/issues`, is kept for the issue files T-2166 writes.
- **The figures fit.** The forecast stays 3h, against flai's 1h12m, with the reason stated. The value is 33.33 USD/week, flai's share of E-0029's 150 USD/week by duration.

**One point for the story's agent, from your first assumption.** Inside the container the guard hook cannot build `bin/flai`, and it passes on any exit but 2, so the guard would be silently off for the agent worked inside. Before T-2166 starts that agent, put the image's `flai` at the worktree's `bin/flai`. It is git-ignored, so this changes no tracked file, and it keeps the guard on for the spike. Check that a refused call is refused before going on, then record what it took in the results. If that cannot be done, stop and ask on a thread rather than run the inner agent unguarded.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md

### 2026-10-26T12:28:20Z orchestrator
Resolved: Plan approved under plan_backlog_stories: the tasks, touches, forecast, and value fit S-0505. The guard-off risk inside the container is noted for T-2166.
