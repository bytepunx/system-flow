---
id: TH-0051
title: flai check --strict fails the close-out on warnings this story does not own
anchor:
  path: wip/kanban/stories/S-0187-the-conventions-say-how-agents-write-shell-commands-for-zsh-and-gate-each-step-of-a-chain.md
  item: S-0187
status: open
participants: [agent-S-0187]
created: 2026-10-01T09:30:24Z
updated: 2026-10-01T09:30:24Z
---

# TH-0051 flai check --strict fails the close-out on warnings this story does not own

On wip/kanban/stories/S-0187-the-conventions-say-how-agents-write-shell-commands-for-zsh-and-gate-each-step-of-a-chain.md.

## Entries

### 2026-10-01T09:30:24Z agent-S-0187
The new `scripts/close-out.sh` runs `flai check --strict`, and in this story's worktree it exits 1 on nine warnings, read from the main checkout's `wip/`. Eight are not S-0187's to fix: E-0003, E-0010, E-0011, E-0012 are done and not archived (`item.archive`); TH-0026 and TH-0032 are answered on archived stories (`threads.archived`); `wip/agents/S-0182.md:31` has an MD029. The ninth is mine: closing I-0006 and I-0012 regenerates `design/issues/summary.md`, which S-0182 (T-0655) also touches, so I declared it and `wip.overlap` warns both ways.

Recommended: you archive the four epics and resolve the two threads (or tell me to: `flai archive` and `flai thread resolve` from this session), S-0182's agent fixes its narrative's list, and the summary.md overlap stands as declared, since the file is regenerated and a conflict there resolves by regenerating it. Then I run the close-out clean and move S-0187 to review. Alternative: I move S-0187 to review with the check's only failures being these, recorded in the narrative. Which do you want?
