---
id: TH-0687
title: "Plan for S-0517: the container mounts the main checkout's .git read-only and the story's registration read-write, after S-0508"
anchor:
  path: wip/kanban/stories/S-0517-a-story-s-worktree-mounted-alone-into-a-container-is-not-a-git-repository-since-its-git-file-names-the-main-checkout-s-git-worktrees-folder.md
  item: S-0517
status: open
participants: [planner-S-0517, orchestrator]
created: 2026-10-28T13:21:29Z
updated: 2026-10-28T13:21:52Z
---

# TH-0687 Plan for S-0517: the container mounts the main checkout's .git read-only and the story's registration read-write, after S-0508

On wip/kanban/stories/S-0517-a-story-s-worktree-mounted-alone-into-a-container-is-not-a-git-repository-since-its-git-file-names-the-main-checkout-s-git-worktrees-folder.md.

## Entries

### 2026-10-28T13:21:29Z planner-S-0517
Plan written by planner-S-0517 on 2026-10-28: tags, topics, touches, after, forecast, cost of delay, and five tasks in four layers.

The remediation assumed, proposed from I-0215's one instance and the code: git inside the container only reads (status, diff, log, show, and the index refresh status writes); every commit is the host's through task_done and the MCP verify, and S-0521 makes flai verify, stream sync, and task done refuse inside the container. So the container mounts two more paths at their host paths: the main checkout's `.git` read-only, and the story's own `.git/worktrees/S-nnnn` read-write over it, where git writes the worktree's index, HEAD, and locks. Objects, refs, packed-refs, config, and hooks stay the host's, as ADR-0027 keeps the dashboard's. A `git commit` in the container is refused by git. The clone of the story branch the issue also names stays the Kubernetes path ADR-0203 keeps open, since S-0505 measured the bind mount as fast enough. An ADR refining ADR-0203 records the choice first.

Tasks and layers:

- Layer 1: T-2218, the ADR.
- Layer 2: T-2219, `storygit.GitMounts` with a real-git test that reproduces I-0215; T-2220, `scripts/agent-container.sh` mounts the two paths by default and `FLAI_AGENT_GIT_MOUNT` goes, with the runbook and scripts index.
- Layer 3: T-2221, the container adapter's `Start` adds the two mounts from the helper and its argv test asserts nothing else under `.git` is writable.
- Layer 4: T-2222, the users guide, agent-adapters.md, and docker.md, a check inside the container, and `flai issue close I-0215`.

Assumptions, each worth a word if wrong:

1. `after: S-0508`. The fix belongs in the container adapter's `Start` (`flai/internal/harness/container.go`), which S-0508 writes, so this story needs S-0508 on main; the issue itself names the remediation as the adapter's. The alternative is to fold the two mounts into S-0508 by amending its criterion 3 ("the three mounts and no other") and leave this story the ADR, the helper, the script, and the issue close with no `after`; say so on this thread and I, or S-0508's agent, adjust.
2. Read-only suffices: the agent in the container never needs `git add`, `commit`, or `stash`, since task_done commits on the host. If the operator wants the agent to commit inside the container, the registration mount is not enough and the ADR must say which refs become writable.
3. The whole `.git` read-only, not the objects and refs one by one: other worktrees' registrations and every branch's objects are then readable from the container, as they are on the host today; mounting pieces is brittle against git's layout.

Figures: forecast raised from flai's 4m (size 2 with no touches) to 1h15m for the ADR, the helper, the adapter, the script, a hand check in the image, and four documents; delivery 2026-10-29T02:41Z from flai's play-out start; cost of delay 37.50 USD a week as `flai cod` gives it from the issue's 15m per cycle. No operator input is missing.

Overlaps: S-0521 (in progress) on the script, the runbook, and `docs/users/flai.md`; S-0414 and S-0502 (in progress) on `docs/users/flai.md`; S-0507 (backlog) on the script. The story is held until they move to review and waits on S-0508 anyway.

### 2026-10-28T13:21:52Z orchestrator (recommendation)
**Recommended for the operator:** approve the plan as written, with T-2218's ADR refining ADR-0203. Keep `after: S-0508` (assumption 1) rather than amending S-0508's criterion 3.

## The plan fits

- **Tasks:** five tasks in four layers. T-2218, the ADR, comes first. T-2219, the helper with a real-git test reproducing I-0215, and T-2220, the script, share no path. T-2221 puts the mounts in the adapter's `Start`, and T-2222 documents them and closes I-0215.
- **Figures:** forecast 1h15m, kept over flai's 4m with the reason stated. Value 37.50 USD a week, from `flai cod` on the issue's 15m input.
- **`after: S-0508`:** right. The adapter's `Start` is S-0508's file, and S-0508's criteria are already finalized.

## Why it is yours

ADR-0203 § Decision names the container's mounts and ends them with "Nothing else: no home, no `~/.ssh`, no `~/.flai`, no main checkout, no `wip/`". It also says the container holds "no git credential and no `.git` of the main checkout". This plan mounts the main checkout's `.git`, so it reverses an accepted decision. Only you can do that.

## Why I recommend it anyway

- **Read-only keeps the guarantee that matters.** The container can commit nothing, push nothing, and change no ref, hook or config. Every commit stays the host's through `task_done`, and S-0521 makes flai refuse the commit paths inside the container.
- **It is what git needs.** Without the common dir, the worktree is not a repository inside the container at all (I-0215). The story's own registration, mounted read-write over it, is all git writes for `status`.
- **What it costs:** every branch's objects, and other worktrees' registrations, become readable from the container, as they already are to an agent on the host. If you want less, the ADR can mount objects and refs piece by piece. The planner warns that this is brittle against git's layout.

**The alternative:** the clone of the story branch that ADR-0203 keeps for Kubernetes. S-0505 measured the bind mount as fast enough, so the clone would add cost for no gain here.

The thread stays open for your answer.

Source: design/adrs/0203-flai-may-run-a-story-s-agent-in-a-container-from-a-baseline-image-it-builds.md § Decision
