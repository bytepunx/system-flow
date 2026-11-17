---
id: T-3551
type: task
nature: improvement
title: An ADR decides what a rebase and resolve job's close-out owes when main moves under it, and how its record names the round
status: backlog
parent: S-0824
owner: alex
created: 2026-11-17T08:33:20Z
updated: 2026-11-17T08:33:20Z
transitions: []
stream: S-0824
tags: [flai, adr]
touches: [design/adrs/README.md]
---
# T-3551 An ADR decides what a rebase and resolve job's close-out owes when main moves under it, and how its record names the round

## Work

Write an ADR with `flai adr new`, refining ADR-0313 and ADR-0156, that decides what a rebase and resolve job's close-out owes when main moves under it. I-0398 records the case: job S-0817-20261117T042837Z ran three rounds and 33 minutes, because the acceptances of S-0742 and S-0753 each brought in 36 or 37 paths that selected every tier.

Weigh the issue's three options:

- **Pass at the head it verified.** In a resolve job, the close-out's last check, `flai verify --sync-only`, passes when main moved but the acceptance's `verify_behind_main` gate would raise no blocker for that head. The gate is ADR-0291 as narrowed by ADR-0303: no path main brought in since the round's merge base overlaps the story's change set, or the overlapping paths send back no tier the report ran. The acceptance judges anything main brings in later, as it does today. Another round runs only when the gate would block. The plan assumes this option, and the plan's thread asks the operator to confirm it.
- **Overlap-based reuse in later rounds.** `sentBack` in `flai/internal/verify/story.go` would select tiers by ADR-0303's overlap in place of every tier main's paths select. This changes what every close-out reuses, so weigh it against ADR-0156 §3.
- **One order for acceptances and resolve jobs.** ADR-0156 rejected holding acceptances because it delays the operator.

Say how ADR-0312's acceptance suites (S-0808) bear on the decision. Say also whether the agent's own close-out, before review, keeps ADR-0156's rule. The plan assumes it does: a story moving to review must contain main.

Decide how the job says which round it is on and why it started another:

- the round line the scripts and the container close-out print, `close-out: <why>; flai stream sync <story> (round N of M)`, with a why that names main's overlapping paths or the tiers they send back;
- the fields `closeout.Status` and the job's record carry;
- what the review page shows.

Name the flag or variable the resolve job passes, for example `CLOSE_OUT_RESOLVE=1` and a `flai verify --sync-only` flag. The later tasks build what this ADR names, so every other task waits for this one.

## Done when

- The ADR is accepted in `design/adrs`, refines ADR-0313 and ADR-0156, and `design/adrs/README.md` lists it.
- It names the condition under which the last check passes over a moved main, the flag or variable that asks for it, the round fields, and the alternatives with why each was taken or rejected.
- `flai check --strict` passes on it.

## Notes
