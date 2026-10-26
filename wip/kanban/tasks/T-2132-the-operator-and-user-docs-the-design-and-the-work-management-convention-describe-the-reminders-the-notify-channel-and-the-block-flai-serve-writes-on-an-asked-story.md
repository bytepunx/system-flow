---
id: T-2132
type: task
nature: improvement
title: The operator and user docs, the design, and the work-management convention describe the reminders, the notify channel, and the block flai serve writes on an asked story
status: in-progress
parent: S-0498
owner: alex
created: 2026-10-26T04:15:16Z
updated: 2026-10-26T04:57:11Z
transitions:
  - to: ready
    at: 2026-10-26T04:57:10Z
    by: agent-S-0498
  - to: in-progress
    at: 2026-10-26T04:57:11Z
    by: agent-S-0498
stream: S-0498
tags: [flai, template]
touches: [docs/operators/index.md, docs/users/flai.md, design/system/workflow.md, design/system/flai-cli.md, design/system/strategic-agents.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md]
after: [T-2129, T-2130]
---
# T-2132 The operator and user docs, the design, and the work-management convention describe the reminders, the notify channel, and the block flai serve writes on an asked story

## Work

Write what the code tasks changed, in the same story, as documentation.md asks:

- `docs/operators/index.md`: a section `## Reminding the operator of what awaits them`, after `## Keeping the host's flai current`: what `threads.remind_after` and `threads.remind_every` do, how to set the channel with `flai serve notify set`, two worked channels for this host (a `curl` to an `ntfy` topic as `--command`; a webhook `--url`), that a reminder writes nothing in the repository, and that a story whose agent ended asking shows blocked with `awaiting the operator on TH-nnnn` until the answer.
- `docs/users/flai.md`: under `### Blocking`, that `flai serve` blocks a story whose agent ended on an open question to you and unblocks it when your answer starts the agent again, so `flai board` shows what waits on you; under `### Threads`, that a thread awaiting you is reminded on the host's channel and that the reminder is not an entry; and a cross-reference from the `flai serve` section to T-2130's `flai serve notify` paragraph.
- `design/system/workflow.md` § Blocking: the block `flai serve` writes, with the ADR linked; § Cadence: the reminder's cadence.
- `design/system/flai-cli.md`: `flai serve notify` and the reminder loop, with the manifest keys and the serve-folder state, in the `flai serve` text.
- `design/system/strategic-agents.md` § Answering threads: a recommendation ends no wait, and `flai serve` reminds the operator of the pending one past `threads.remind_after`; one sentence with the ADR linked.
- `design/conventions/work-management.md`, in the baseline rule that starts "When you need the designer to decide something": when you end your session on the question, `flai serve` blocks the story with the thread as the reason and unblocks it when the answer starts you again, so you need not run `flai block` for it; the same words in `template/root/design/conventions/work-management.md`, and a `template/CHANGELOG.md` entry.

Waits for T-2129 and T-2130, whose behaviour it describes; T-2130 also changes `docs/users/flai.md`.

## Done when

- Each document above names the two manifest keys, `flai serve notify`, or the block, as its audience needs, and links the ADR of T-2127.
- The convention's baseline and its template copy say the same words, and `flai check --strict` and the markdown lint pass on every changed file.

## Notes

Layer 4, with T-2131. The reminder's settings rows are T-2128's and T-2130's, not this task's.
