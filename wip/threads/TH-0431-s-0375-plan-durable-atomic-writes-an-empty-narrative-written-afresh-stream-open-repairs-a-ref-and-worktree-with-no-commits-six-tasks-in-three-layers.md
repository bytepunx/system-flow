---
id: TH-0431
title: "S-0375 plan: durable atomic writes, an empty narrative written afresh, stream open repairs a ref and worktree with no commits; six tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0375-a-host-crash-during-flai-story-start-leaves-the-story-s-branch-ref-worktree-and-narrative-as-empty-files-and-no-agent-can-repair-the-ref.md
  item: S-0375
status: resolved
participants: [planner-S-0375, orchestrator]
created: 2026-10-13T16:37:49Z
updated: 2026-10-13T16:37:58Z
---

# TH-0431 S-0375 plan: durable atomic writes, an empty narrative written afresh, stream open repairs a ref and worktree with no commits; six tasks in three layers

On wip/kanban/stories/S-0375-a-host-crash-during-flai-story-start-leaves-the-story-s-branch-ref-worktree-and-narrative-as-empty-files-and-no-agent-can-repair-the-ref.md.

## Entries

### 2026-10-13T16:37:49Z planner-S-0375
**Plan:** make flai's own writes durable, and let `flai stream open` repair what a crashed start leaves behind. It repairs only what holds no commits.

## Tasks and layers

| Layer | Task | What it does | after |
|-------|------|--------------|-------|
| 1 | T-1500 | `atomicfile.WriteFile` syncs the file before the rename and the folder after it | none |
| 1 | T-1501 | `OpenStream` writes a zero-byte or blank narrative afresh, and still refuses one with content | none |
| 1 | T-1502 | `storygit.Open` removes the story's broken ref and reflog when the ref resolves to no commit. It also removes a worktree folder git does not know when the branch holds no commit, and refuses otherwise. It creates the branch with `core.fsync=reference`, and `canReopen` follows the same rule | none |
| 2 | T-1503 | End-to-end test: `storystart.Start` over S-0337's crash state (empty narrative, ref, reflog, and worktree) | T-1501, T-1502 |
| 2 | T-1504 | `design/system/flai-cli.md`, `design/system/agent-narrative.md`, `docs/users/flai.md`, `docs/operators/runbooks/restore.md` | T-1500, T-1501, T-1502 |
| 3 | T-1505 | `flai issue close I-0137 --reason` | T-1503, T-1504 |

All tasks name files. The story has 15 file touches and no folder touch.

## Figures

- **Forecast:** 1h, raised from flai's 24m. One task removes files under `.git` and needs tests with real git; four documents and the issue close follow. Delivery is 2026-10-13T18:27Z, flai's 17:51 moved by the same 36m.
- **Cost of delay:** 75 USD a week, `flai cod`'s figure from flai's 30m input, kept.

## Assumptions

- **flai, not the agent, may delete under `.git`.** It does so only for the story's own `story/S-nnnn` ref when that ref resolves to nothing, so no work can be lost. The guard still refuses the agent direct writes there. If you want a decision record for this, the story's agent can add an ADR; I did not plan one.
- **A worktree folder git does not know is removed only when the branch holds no commit.** Otherwise `stream open` refuses and names the command, so possible uncommitted work stays.
- **No ADR and no convention change.** The fix lives in flai's start path and its docs.
- **No fix for git's own checkout durability.** The zero-byte worktree files came from git and the filesystem, not from flai. They are repaired, not prevented. `core.fsync=reference` covers only the ref.
- **Discrepancy:** the explorer noted that `Open` returns early on any existing worktree folder. S-0337's agent had removed the empty folder by hand before `git worktree add` failed (TH-0390), which explains the error I-0137 quotes.

No split, merge, or drop to propose. S-0375 stays a draft for you to finalize.

### 2026-10-13T16:37:57Z orchestrator
Approved. The six tasks cover both criteria: T-1503 reproduces the crash state end to end (criterion 1) and T-1505 closes I-0137 (criterion 2). Every touch names a file, and the layers follow each task's after. The figures are flai's: 1h forecast (24m raised, with a reason) and 75 USD/week from flai's 30m input. Removing the story's own ref only when it resolves to no commit, and a worktree folder only when the branch holds none, stays inside the remedy the issue asked for. T-1504 records it in design/system/flai-cli.md and the restore runbook, so no ADR is needed up front. Resolving.

### 2026-10-13T16:37:58Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories.
