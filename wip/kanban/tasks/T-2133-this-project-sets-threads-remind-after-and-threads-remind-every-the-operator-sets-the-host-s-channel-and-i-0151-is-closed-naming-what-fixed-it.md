---
id: T-2133
type: task
nature: improvement
title: This project sets threads.remind_after and threads.remind_every, the operator sets the host's channel, and I-0151 is closed naming what fixed it
status: backlog
parent: S-0498
owner: alex
created: 2026-10-26T04:15:49Z
updated: 2026-10-26T04:15:49Z
transitions: []
stream: S-0498
tags: [flai]
touches: [system-flow.yaml, design/issues/I-0151-a-pending-recommendation-waits-a-day-or-more-for-the-operator-s-confirmation-while-its-story-idles-in-progress.md, design/issues/summary.md]
after: [T-2131, T-2132]
---
# T-2133 This project sets threads.remind_after and threads.remind_every, the operator sets the host's channel, and I-0151 is closed naming what fixed it

## Work

Make the remedy hold here and close the issue:

- In the story's worktree, `flai manifest set threads.remind_after 4h` and `flai manifest set threads.remind_every 24h`, so this project reminds the operator of a thread that has awaited them half a working day, and daily after that, from the release that carries this story. The values are the planner's assumption: the operator may name others on the plan's thread.
- The channel needs the operator's hands, since only they may write `~/.flai/config.json` (safety.md): open one thread on S-0498 giving the exact `flai serve notify set` command for this host, with a `--command` that reaches them off the dashboard (an `ntfy` topic through `curl`, or a Windows toast from WSL2 through `powershell.exe`), and ask them to run it once the story's release is installed. The thread does not hold the story: the setting is theirs to make after acceptance.
- Run the story's tests once more, then close the issue: `flai issue close I-0151 --reason "<what fixed it>"`, naming the ADR, the reminders (`threads.remind_after`, `flai serve notify`), the block on an asked story, and S-0502 for the design-choice waits it leaves to `confirm_after`. The close regenerates `design/issues/summary.md`.
- Tick criterion 1 once T-2129's and T-2131's tests pass and criterion 2 once the issue is closed, with `flai criteria tick`.

Waits for T-2131 and T-2132, so the issue closes on what is built and documented.

## Done when

- `system-flow.yaml` carries `threads.remind_after: 4h` and `threads.remind_every: 24h`, and `flai check --strict` passes on it.
- I-0151 is closed with a reason naming the ADR, the reminder, the block, and S-0502, `design/issues/summary.md` no longer lists it as open, and both criteria of S-0498 are ticked.
- A thread on S-0498 gives the operator the `flai serve notify set` command for this host.

## Notes

Layer 5, the last. The manifest values 4h and 24h are assumptions recorded on the plan's thread.
