---
id: TH-0027
title: "Commit outstanding work: start an agent, or have flai commit it itself?"
anchor:
  path: wip/kanban/stories/S-0140-agents-should-not-leave-uncommitted-work-in-their-branch-when-moving-to-review.md
  item: S-0140
status: resolved
participants: [agent-S-0140, alex]
created: 2026-09-26T21:11:01Z
updated: 2026-09-27T03:55:44Z
---

# TH-0027 Commit outstanding work: start an agent, or have flai commit it itself?

On wip/kanban/stories/S-0140-agents-should-not-leave-uncommitted-work-in-their-branch-when-moving-to-review.md.

## Entries

### 2026-09-26T21:11:01Z agent-S-0140
The second criterion asks for a serve action that has outstanding work in a story's worktree committed. Two ways:

**(a) Recommended: start an agent to do it.** A new host method (`agent.commit`), `flai serve agent commit S-nnnn`, and a "Have an agent commit it" button on the story page. It shows when the acceptance preview finds uncommitted changes in the story worktree. It starts the story's agent the way a restart does. The prompt says: commit everything outstanding in the worktree with a proper message, run lint, do not change anything else, and leave the story in review. The agent can judge whether a file is junk to discard or work to commit. It needs the agent action enabled and a harness.

**(b) flai commits it itself.** `git add -A` and `git commit -m "chore: [S-nnnn] commit outstanding work"` in the worktree, from the same button. It is instant and needs no agent. But it commits blindly, with a generic message.

For criterion 1, `flai move S-nnnn review` (CLI and MCP `item_move`) will refuse a story whose worktree has uncommitted changes, naming the files. The prompt flai serve gives agents will also say to commit everything before moving to review.

I am building (a) unless you say otherwise.

### 2026-09-27T03:55:32Z alex
A

### 2026-09-27T03:55:44Z agent-S-0140
Resolved: The designer chose (a): an agent is started to commit the worktree, as built.
