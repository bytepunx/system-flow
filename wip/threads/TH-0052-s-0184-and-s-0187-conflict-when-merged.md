---
id: TH-0052
title: S-0184 and S-0187 conflict when merged
anchor:
  path: wip/kanban/stories/S-0187-the-conventions-say-how-agents-write-shell-commands-for-zsh-and-gate-each-step-of-a-chain.md
  item: S-0187
status: answered
participants: [flai, agent-S-0187]
created: 2026-10-01T09:30:36Z
updated: 2026-10-01T09:30:53Z
---

# TH-0052 S-0184 and S-0187 conflict when merged

On wip/kanban/stories/S-0187-the-conventions-say-how-agents-write-shell-commands-for-zsh-and-gate-each-step-of-a-chain.md.

## Entries

### 2026-10-01T09:30:36Z flai
A trial merge of story/S-0184 with story/S-0187 at flai stream sync conflicts in:

- `design/issues/summary.md`

Whichever of S-0184 and S-0187 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-01T09:30:53Z agent-S-0187
Settled: S-0187 now names S-0184 (and S-0186, TH-0053) in `after:`. The conflict is in `design/issues/summary.md`, which `flai issue` regenerates; S-0187 closes I-0006 and I-0012 there. S-0184 is in review and goes first. When S-0187 rebases onto it, the resolution is to take main's file and run `flai issue summary` again.
